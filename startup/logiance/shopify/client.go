package shopify

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shridhar2104/logilo/shopify/pb"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"google.golang.org/grpc"
)

// OrdersResponse represents paginated orders response
type OrdersResponse struct {
	Orders      []*Order
	TotalCount  int32
	TotalPages  int32
	CurrentPage int32
}

// Client struct for gRPC communication.
type Client struct {
	conn    *grpc.ClientConn
	service pb.ShopifyServiceClient
}

// NewClient creates a new gRPC client.
func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewShopifyServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    req := &pb.GetOrderRequest{
        OrderId: orderID,
    }

    resp, err := c.service.GetOrder(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }

    // Handle case where order is not found
    if resp == nil {
        return nil, fmt.Errorf("order not found")
    }

    // Convert string ID to int64
    id, err := strconv.ParseInt(resp.Id, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("failed to parse order ID: %w", err)
    }

    createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
    if err != nil {
        createdAt = time.Time{}
    }

    return &Order{
        ID:                id,
        Name:             resp.OrderName,
        CreatedAt:        createdAt,
        Currency:         resp.Currency,
        TotalPrice:       float64(resp.TotalPrice),
        SubtotalPrice:    float64(resp.SubtotalPrice),
        TotalTax:         float64(resp.TotalTax),
        FinancialStatus:  resp.FinancialStatus,
        FulfillmentStatus: resp.FulfillmentStatus,
        Customer: goshopify.Customer{
            Email: resp.CustomerEmail,
            FirstName: strings.Split(resp.CustomerName, " ")[0],
            LastName: strings.Join(strings.Split(resp.CustomerName, " ")[1:], " "),
        },
    }, nil
}

// GetOrdersForAccount fetches orders for all shops under an account with pagination
func (c *Client) GetOrdersForAccount(ctx context.Context, accountId string, page, pageSize int32) (*OrdersResponse, error) {
	req := &pb.GetOrdersForAccountRequest{
		AccountId: accountId,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.service.GetOrdersForAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]*Order, len(resp.Orders))
	for i, o := range resp.Orders {
		createdAt, err := time.Parse(time.RFC3339, o.CreatedAt)
		if err != nil {
			createdAt = time.Time{}
		}

		// Convert string ID to int64
		id, err := strconv.ParseInt(o.Id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order ID: %w", err)
		}

		// Create order with available fields
		orders[i] = &Order{
			ID:                id,
			Name:             o.OrderName,
			CreatedAt:        createdAt,
			Currency:         o.Currency,
			TotalPrice:       float64(o.TotalPrice),
			SubtotalPrice:    float64(o.SubtotalPrice),
			TotalTax:         float64(o.TotalTax),
			FinancialStatus:  o.FinancialStatus,
			FulfillmentStatus: o.FulfillmentStatus,
			Customer: goshopify.Customer{
				Email: o.CustomerEmail,
			},
		}
	}

	return &OrdersResponse{
		Orders:      orders,
		TotalCount:  resp.TotalCount,
		TotalPages:  resp.TotalPages,
		CurrentPage: resp.CurrentPage,
	}, nil
}



// // GetOrdersForShopAndAccount fetches orders for a specific shop and account.
// func (c *Client) GetOrdersForShopAndAccount(ctx context.Context, shopName, accountId string) ([]*Order, error) {
// 	res, err := c.service.GetOrdersForShopAndAccount(ctx, &pb.GetOrdersForShopAndAccountRequest{
// 		ShopName:  shopName,
// 		AccountId: accountId,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	orders := make([]*Order, len(res.Orders))
// 	for i, o := range res.Orders {
// 		orders[i] = &Order{
// 			ID:         o.Id,
// 			TotalPrice: float64(o.TotalPrice),
// 		}
// 	}
// 	return orders, nil
// }
// //http://djcajdjd?code=jfl&shopname

// GenerateAuthURL generates an authorization URL for a Shopify store.
func (c *Client) GenerateAuthURL(ctx context.Context, shopName string) (string, error) {
	req := &pb.GetAuthorizationURLRequest{
		ShopName:  shopName,
		State: "your_unique_nonce",
		
	}
	resp, err := c.service.GetAuthorizationURL(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate auth URL: %v", err)
	}
	return resp.AuthUrl, nil
}

// ExchangeAccessToken exchanges a Shopify auth code for an access token.
func (c *Client) ExchangeAccessToken(ctx context.Context, shopName, code, accountId string) error {
	req := &pb.ExchangeAccessTokenRequest{
		ShopName:  shopName,
		Code:      code,
		AccountId: accountId,
	}
	_, err := c.service.ExchangeAccessToken(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to exchange access token: %v", err)
	}
	return nil
}


// client.go
func (c *Client) SyncOrders(ctx context.Context, accountId string) (map[string]*pb.ShopSyncStatus, error) {
    req := &pb.SyncOrderRequest{
        AccountId: accountId,
    }
    
    resp, err := c.service.SyncOrders(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to sync orders: %v", err)
    }
    
    return resp.ShopResults, nil
}