package shopify

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shridhar2104/logilo/shopify/pb"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"
)

const (
	batchSize          = 250
	initialRateLimit   = 200 * time.Millisecond // 5 requests per second
	maxRateLimit      = 1000 * time.Millisecond // 1 request per second
	minRateLimit      = 100 * time.Millisecond  // 10 requests per second
	maxRetries        = 5
	orderSyncTimeout  = 30 * time.Minute
	batchSyncTimeout  = 5 * time.Minute
)

type Service interface {
	GenerateAuthURL(ctx context.Context, shopName, state string) (string, error)
	ExchangeAccessToken(ctx context.Context, shop, code, accountId string) error
	SyncOrdersForAccount(ctx context.Context, accountId string) (map[string]*pb.ShopSyncStatus, error)
	GetAccountOrders(ctx context.Context, accountId string, page, pageSize int) ([]Order, int, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)

}

type ShopifyService struct {
	App  goshopify.App
	Repo Repository
}

type adaptiveRateLimiter struct {
	ticker    *time.Ticker
	interval  time.Duration
	mu        sync.Mutex
}

func newAdaptiveRateLimiter() *adaptiveRateLimiter {
	return &adaptiveRateLimiter{
		ticker:   time.NewTicker(initialRateLimit),
		interval: initialRateLimit,
	}
}

func (r *adaptiveRateLimiter) adjustRate(isRateLimited bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if isRateLimited {
		// Slow down
		r.interval = min(r.interval*2, maxRateLimit)
	} else {
		// Speed up gradually
		r.interval = max(r.interval-50*time.Millisecond, minRateLimit)
	}
	r.ticker.Reset(r.interval)
}

func (r *adaptiveRateLimiter) wait() {
	<-r.ticker.C
}

func (r *adaptiveRateLimiter) stop() {
	r.ticker.Stop()
}

type syncResult struct {
	shopName     string
	success      bool
	errorMessage string
	ordersSynced int32
}

func NewShopifyService(apiKey, apiSecret, redirectURL string, repo Repository) Service {
	return &ShopifyService{
		App: goshopify.App{
			ApiKey:      apiKey,
			ApiSecret:   apiSecret,
			RedirectUrl: redirectURL,
			Scope:       "read_products,read_orders,read_fulfillments,read_all_orders,read_merchant_managed_fulfillment_orders,write_merchant_managed_fulfillment_orders,write_fulfillments",
		},
		Repo: repo,
	}
}

// GetOrder retrieves a single order by ID
func (s *ShopifyService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    if s == nil || s.Repo == nil {
        return nil, fmt.Errorf("shopify service not properly initialized")
    }

    if orderID == "" {
        return nil, fmt.Errorf("order ID cannot be empty")
    }

    order, err := s.Repo.GetOrder(ctx, orderID)
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }

    return order, nil
}


// Implementation for getting all orders for an account
func (s *ShopifyService) GetAccountOrders(ctx context.Context, accountId string, page, pageSize int) ([]Order, int, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 50 // default page size
    }

    offset := (page - 1) * pageSize
    return s.Repo.GetAccountOrders(ctx, accountId, pageSize, offset)
}


func (s *ShopifyService) SyncOrdersForAccount(ctx context.Context, accountId string) (map[string]*pb.ShopSyncStatus, error) {
	if s == nil || s.Repo == nil {
		return nil, fmt.Errorf("shopify service not properly initialized")
	}

	shopTokens, err := s.Repo.GetShopTokens(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop tokens: %w", err)
	}

	if len(shopTokens) == 0 {
		return map[string]*pb.ShopSyncStatus{}, nil
	}

	results := make(chan syncResult, len(shopTokens))
	var wg sync.WaitGroup

	for shopName, token := range shopTokens {
		if shopName == "" || token == "" {
			results <- syncResult{
				shopName:     shopName,
				success:      false,
				errorMessage: "invalid shop name or token",
				ordersSynced: 0,
			}
			continue
		}

		wg.Add(1)
		go func(shopName, token string) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					results <- syncResult{
						shopName:     shopName,
						success:      false,
						errorMessage: fmt.Sprintf("panic recovered: %v", r),
						ordersSynced: 0,
					}
				}
			}()

			s.syncShopOrders(ctx, shopName, token, accountId, results)
		}(shopName, token)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	syncResults := make(map[string]*pb.ShopSyncStatus)
	for result := range results {
		if result.shopName != "" {
			syncResults[result.shopName] = &pb.ShopSyncStatus{
				Success:      result.success,
				ErrorMessage: result.errorMessage,
				OrdersSynced: result.ordersSynced,
			}
		}
	}

	return syncResults, nil
}

func (s *ShopifyService) syncShopOrders(ctx context.Context, shopName, token, accountId string, results chan<- syncResult) {
    client, err := goshopify.NewClient(s.App, shopName, token)
    if err != nil {
        results <- syncResult{
            shopName:     shopName,
            success:      false,
            errorMessage: fmt.Sprintf("failed to create client: %v", err),
            ordersSynced: 0,
        }
        return
    }

    orderCtx, cancel := context.WithTimeout(ctx, orderSyncTimeout)
    defer cancel()

    // Get last sync timestamp for incremental sync
    lastSync, err := s.Repo.GetLatestOrderTimestamp(orderCtx, shopName, accountId)
    if err != nil {
        // Log error but continue with full sync
        fmt.Printf("Failed to get latest order timestamp for %s: %v\n", shopName, err)
    }

    // Create order list options
    options := goshopify.OrderListOptions{
        ListOptions: goshopify.ListOptions{
            Limit: batchSize,
        },
        Status: "any",
    }

    if !lastSync.IsZero() {
        options.UpdatedAtMin = lastSync
    }

    rateLimiter := newAdaptiveRateLimiter()
    defer rateLimiter.stop()

    var ordersSynced int32
    var currentBatch []Order
    currentBatch = make([]Order, 0, batchSize)

    backoff := time.Second
    retryCount := 0

    for {
        select {
        case <-orderCtx.Done():
            results <- syncResult{
                shopName:     shopName,
                success:      false,
                errorMessage: fmt.Sprintf("operation cancelled or timed out: %v", orderCtx.Err()),
                ordersSynced: ordersSynced,
            }
            return
        default:
            rateLimiter.wait()

            shopifyOrders, pagination, err := client.Order.ListWithPagination(orderCtx, &options)
            if err != nil {
                if strings.Contains(err.Error(), "Exceeded") && retryCount < maxRetries {
                    retryCount++
                    rateLimiter.adjustRate(true)
                    time.Sleep(backoff)
                    backoff *= 2
                    if backoff > maxRateLimit {
                        backoff = maxRateLimit
                    }
                    continue
                }

                results <- syncResult{
                    shopName:     shopName,
                    success:      false,
                    errorMessage: fmt.Sprintf("failed to list orders after %d retries: %v", retryCount, err),
                    ordersSynced: ordersSynced,
                }
                return
            }

            rateLimiter.adjustRate(false)
            backoff = time.Second
            retryCount = 0

            // Process current batch
            for _, shopifyOrder := range shopifyOrders {
                order := convertShopifyOrderToOrder(&shopifyOrder)
                currentBatch = append(currentBatch, order)

                if len(currentBatch) >= batchSize {
                    if err := s.syncBatch(ctx, currentBatch, shopName, accountId); err != nil {
                        results <- syncResult{
                            shopName:     shopName,
                            success:      false,
                            errorMessage: fmt.Sprintf("failed to sync batch: %v", err),
                            ordersSynced: ordersSynced,
                        }
                        return
                    }
                    ordersSynced += int32(len(currentBatch))
                    currentBatch = currentBatch[:0]
                }
            }

            // Handle remaining orders in the last batch
            if len(currentBatch) > 0 && (pagination == nil || pagination.NextPageOptions == nil) {
                if err := s.syncBatch(ctx, currentBatch, shopName, accountId); err != nil {
                    results <- syncResult{
                        shopName:     shopName,
                        success:      false,
                        errorMessage: fmt.Sprintf("failed to sync final batch: %v", err),
                        ordersSynced: ordersSynced,
                    }
                    return
                }
                ordersSynced += int32(len(currentBatch))
            }

            if pagination == nil || pagination.NextPageOptions == nil {
                results <- syncResult{
                    shopName:     shopName,
                    success:      true,
                    errorMessage: "",
                    ordersSynced: ordersSynced,
                }
                return
            }

            options.ListOptions.Page = pagination.NextPageOptions.Page
        }
    }
}
func (s *ShopifyService) syncBatch(ctx context.Context, batch []Order, shopName, accountId string) error {
	batchCtx, cancel := context.WithTimeout(ctx, batchSyncTimeout)
	defer cancel()

	success, err := s.Repo.SyncOrders(batchCtx, batch, shopName, accountId)
	if err != nil {
		return fmt.Errorf("failed to sync batch: %w", err)
	}
	if !success {
		return fmt.Errorf("batch sync reported failure")
	}
	return nil
}

func (s *ShopifyService) GenerateAuthURL(ctx context.Context, shopName, state string) (string, error) {
	return s.App.AuthorizeUrl(shopName, state)
}

func (s *ShopifyService) ExchangeAccessToken(ctx context.Context, shop, code, accountId string) error {
	accessToken, err := s.App.GetAccessToken(ctx, shop, code)
	if err != nil {
		return err
	}
	return s.Repo.SaveShopCredentials(ctx, shop, accountId, accessToken)
}

func convertShopifyOrderToOrder(shopifyOrder *goshopify.Order) Order {
	order := Order{
		ID:                int64(shopifyOrder.Id),
		Name:             shopifyOrder.Name,
		Email:            shopifyOrder.Email,
		CreatedAt:        *shopifyOrder.CreatedAt,
		UpdatedAt:        *shopifyOrder.UpdatedAt,
		CancelledAt:      shopifyOrder.CancelledAt,
		ClosedAt:         shopifyOrder.ClosedAt,
		ProcessedAt:      shopifyOrder.ProcessedAt,
		Currency:         shopifyOrder.Currency,
		TotalPrice:       convertDecimalToFloat(shopifyOrder.TotalPrice),
		SubtotalPrice:    convertDecimalToFloat(shopifyOrder.SubtotalPrice),
		TotalDiscounts:   convertDecimalToFloat(shopifyOrder.TotalDiscounts),
		TotalTax:         convertDecimalToFloat(shopifyOrder.TotalTax),
		TaxesIncluded:    shopifyOrder.TaxesIncluded,
		FinancialStatus:  string(shopifyOrder.FinancialStatus),
		FulfillmentStatus: string(shopifyOrder.FulfillmentStatus),
		OrderNumber:      shopifyOrder.OrderNumber,
		Test:             shopifyOrder.Test,
		BrowserIp:        shopifyOrder.BrowserIp,
		CancelReason:     string(shopifyOrder.CancelReason),
		Tags:             shopifyOrder.Tags,
		Gateway:          shopifyOrder.Gateway,
		Confirmed:        shopifyOrder.Confirmed,
		ContactEmail:     shopifyOrder.ContactEmail,
		Phone:            shopifyOrder.Phone,
	}

	if shopifyOrder.Customer != nil {
		order.Customer = *shopifyOrder.Customer
	}

	return order
}

func convertDecimalToFloat(d *decimal.Decimal) float64 {
	if d == nil {
		return 0.0
	}
	f, _ := d.Float64()
	return f
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}