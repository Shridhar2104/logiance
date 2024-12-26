package shopify

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/lib/pq"
)
type Repository interface {
    GetShopTokens(ctx context.Context, accountId string) (map[string]string, error)
    SaveShopCredentials(ctx context.Context, shop, accountId, token string) error
    SyncOrders(ctx context.Context, orders []Order, shopName, accountId string) (bool, error)
    GetLatestOrderTimestamp(ctx context.Context, shopName, accountId string) (time.Time, error)
	GetAccountOrders(ctx context.Context, accountId string, limit, offset int) ([]Order, int, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
    Close() 
}
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}
// Close releases the database connection resources.
func (r *PostgresRepository) Close() {
	r.db.Close()
}

// Update in repository.go

func (r *PostgresRepository) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    query := `
        SELECT 
            id, name, email, created_at, updated_at, cancelled_at, 
            closed_at, processed_at, currency, total_price, subtotal_price,
            total_discounts, total_tax, taxes_included, financial_status,
            fulfillment_status, order_number, test, browser_ip,
            cancel_reason, tags, gateway, confirmed, contact_email, phone,
            customer_id, customer_email, customer_first_name, 
            customer_last_name, customer_phone
        FROM orders 
        WHERE id = $1`

    row := r.db.QueryRowContext(ctx, query, orderID)

    var order Order
    var customerID sql.NullInt64
    var customerEmail, customerFirstName, customerLastName, customerPhone sql.NullString
    var cancelledAt, closedAt, processedAt sql.NullTime

    err := row.Scan(
        &order.ID, &order.Name, &order.Email, &order.CreatedAt, &order.UpdatedAt,
        &cancelledAt, &closedAt, &processedAt, &order.Currency,
        &order.TotalPrice, &order.SubtotalPrice, &order.TotalDiscounts,
        &order.TotalTax, &order.TaxesIncluded, &order.FinancialStatus,
        &order.FulfillmentStatus, &order.OrderNumber, &order.Test,
        &order.BrowserIp, &order.CancelReason, &order.Tags, &order.Gateway,
        &order.Confirmed, &order.ContactEmail, &order.Phone,
        &customerID, &customerEmail, &customerFirstName, 
        &customerLastName, &customerPhone,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("order not found: %s", orderID)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to scan order row: %w", err)
    }

    // Handle nullable timestamps
    if cancelledAt.Valid {
        order.CancelledAt = &cancelledAt.Time
    }
    if closedAt.Valid {
        order.ClosedAt = &closedAt.Time
    }
    if processedAt.Valid {
        order.ProcessedAt = &processedAt.Time
    }

    // Handle customer information
    if customerID.Valid {
        order.Customer = goshopify.Customer{
            Id:        uint64(customerID.Int64),
            Email:     customerEmail.String,
            FirstName: customerFirstName.String,
            LastName:  customerLastName.String,
            Phone:     customerPhone.String,
        }
    }

    return &order, nil
}
// Implementation for getting all orders for an account
func (r *PostgresRepository) GetAccountOrders(ctx context.Context, accountId string, limit, offset int) ([]Order, int, error) {
    // First, get total count for all shops under this account
    countQuery := `
        SELECT COUNT(*)
        FROM orders
        WHERE account_id = $1`
    
    var totalCount int
    err := r.db.QueryRowContext(ctx, countQuery, accountId).Scan(&totalCount)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get total order count: %w", err)
    }

    // Then get orders with pagination
    query := `
        SELECT 
            id, name, email, created_at, updated_at, cancelled_at, 
            closed_at, processed_at, currency, total_price, subtotal_price,
            total_discounts, total_tax, taxes_included, financial_status,
            fulfillment_status, order_number, test, browser_ip,
            cancel_reason, tags, gateway, confirmed, contact_email, phone,
            shop_name, customer_id, customer_email, customer_first_name, 
            customer_last_name, customer_phone
        FROM orders 
        WHERE account_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

    rows, err := r.db.QueryContext(ctx, query, accountId, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query orders: %w", err)
    }
    defer rows.Close()

    var orders []Order
    for rows.Next() {
        var order Order
        var customerID sql.NullInt64
        var customerEmail, customerFirstName, customerLastName, customerPhone sql.NullString
        var cancelledAt, closedAt, processedAt sql.NullTime
        var shopName string

        err := rows.Scan(
            &order.ID, &order.Name, &order.Email, &order.CreatedAt, &order.UpdatedAt,
            &cancelledAt, &closedAt, &processedAt, &order.Currency,
            &order.TotalPrice, &order.SubtotalPrice, &order.TotalDiscounts,
            &order.TotalTax, &order.TaxesIncluded, &order.FinancialStatus,
            &order.FulfillmentStatus, &order.OrderNumber, &order.Test,
            &order.BrowserIp, &order.CancelReason, &order.Tags, &order.Gateway,
            &order.Confirmed, &order.ContactEmail, &order.Phone,
            &shopName, &customerID, &customerEmail, &customerFirstName, &customerLastName,
            &customerPhone,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan order row: %w", err)
        }

        // Handle nullable timestamps
        if cancelledAt.Valid {
            order.CancelledAt = &cancelledAt.Time
        }
        if closedAt.Valid {
            order.ClosedAt = &closedAt.Time
        }
        if processedAt.Valid {
            order.ProcessedAt = &processedAt.Time
        }

        // Handle customer information
        if customerID.Valid {
            order.Customer = goshopify.Customer{
                Id:        uint64(customerID.Int64),
                Email:     customerEmail.String,
                FirstName: customerFirstName.String,
                LastName:  customerLastName.String,
                Phone:     customerPhone.String,
            }
        }


        orders = append(orders, order)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("error iterating order rows: %w", err)
    }

    return orders, totalCount, nil
}

func (r *PostgresRepository) GetShopTokens(ctx context.Context, accountId string) (map[string]string, error) {
	query := `
		SELECT shop_name, token 
		FROM tokens 
		WHERE account_id = $1`

	rows, err := r.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to query shop tokens: %w", err)
	}
	defer rows.Close()

	tokens := make(map[string]string)
	for rows.Next() {
		var shopName, token string
		if err := rows.Scan(&shopName, &token); err != nil {
			return nil, fmt.Errorf("failed to scan token row: %w", err)
		}
		tokens[shopName] = token
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating token rows: %w", err)
	}

	return tokens, nil
}

func (r *PostgresRepository) SaveShopCredentials(ctx context.Context, shop, accountId, token string) error {
	query := `
		INSERT INTO tokens (shop_name, account_id, token) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (shop_name, account_id) 
		DO UPDATE SET token = EXCLUDED.token`

	_, err := r.db.ExecContext(ctx, query, shop, accountId, token)
	if err != nil {
		return fmt.Errorf("failed to save shop credentials: %w", err)
	}

	return nil
}

func (r *PostgresRepository) SyncOrders(ctx context.Context, orders []Order, shopName, accountId string) (bool, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the UPSERT statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO orders (
			id, name, email, created_at, updated_at, cancelled_at, 
			closed_at, processed_at, currency, total_price, subtotal_price,
			total_discounts, total_tax, taxes_included, financial_status,
			fulfillment_status, order_number, test, browser_ip,
			cancel_reason, tags, gateway, confirmed, contact_email, phone,
			shop_name, account_id, customer_id, customer_email,
			customer_first_name, customer_last_name, customer_phone
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30, $31, $32
		)
		ON CONFLICT (id) DO UPDATE SET
			updated_at = EXCLUDED.updated_at,
			cancelled_at = EXCLUDED.cancelled_at,
			closed_at = EXCLUDED.closed_at,
			processed_at = EXCLUDED.processed_at,
			total_price = EXCLUDED.total_price,
			subtotal_price = EXCLUDED.subtotal_price,
			total_discounts = EXCLUDED.total_discounts,
			total_tax = EXCLUDED.total_tax,
			financial_status = EXCLUDED.financial_status,
			fulfillment_status = EXCLUDED.fulfillment_status,
			tags = EXCLUDED.tags,
			gateway = EXCLUDED.gateway
	`)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Process orders in smaller batches to avoid memory issues
	batchSize := 100
	for i := 0; i < len(orders); i += batchSize {
		end := i + batchSize
		if end > len(orders) {
			end = len(orders)
		}

		batch := orders[i:end]
		if err := r.processBatch(ctx, stmt, batch, shopName, accountId); err != nil {
			return false, fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return true, nil
}

func (r *PostgresRepository) processBatch(ctx context.Context, stmt *sql.Stmt, orders []Order, shopName, accountId string) error {
	for _, order := range orders {
		customerID := sql.NullInt64{}
		customerEmail := sql.NullString{}
		customerFirstName := sql.NullString{}
		customerLastName := sql.NullString{}
		customerPhone := sql.NullString{}

		if order.Customer.Id != 0 {
			customerID.Int64 = int64(order.Customer.Id)
			customerID.Valid = true
			customerEmail.String = order.Customer.Email
			customerEmail.Valid = true
			customerFirstName.String = order.Customer.FirstName
			customerFirstName.Valid = true
			customerLastName.String = order.Customer.LastName
			customerLastName.Valid = true
			customerPhone.String = order.Customer.Phone
			customerPhone.Valid = true
		}

		_, err := stmt.ExecContext(ctx,
			order.ID, order.Name, order.Email, order.CreatedAt, order.UpdatedAt,
			nullTimePtr(order.CancelledAt), nullTimePtr(order.ClosedAt),
			nullTimePtr(order.ProcessedAt), order.Currency, order.TotalPrice,
			order.SubtotalPrice, order.TotalDiscounts, order.TotalTax,
			order.TaxesIncluded, order.FinancialStatus, order.FulfillmentStatus,
			order.OrderNumber, order.Test, order.BrowserIp, order.CancelReason,
			order.Tags, order.Gateway, order.Confirmed, order.ContactEmail,
			order.Phone, shopName, accountId, customerID, customerEmail,
			customerFirstName, customerLastName, customerPhone,
		)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				// Handle specific PostgreSQL errors
				switch pqErr.Code {
				case "23505": // unique_violation
					continue // Skip duplicate orders
				default:
					return fmt.Errorf("database error: %w", err)
				}
			}
			return fmt.Errorf("failed to insert order %d: %w", order.ID, err)
		}
	}
	return nil
}

func (r *PostgresRepository) GetLatestOrderTimestamp(ctx context.Context, shopName, accountId string) (time.Time, error) {
	query := `
		SELECT COALESCE(MAX(updated_at), '1970-01-01'::timestamp)
		FROM orders
		WHERE shop_name = $1 AND account_id = $2`

	var timestamp time.Time
	err := r.db.QueryRowContext(ctx, query, shopName, accountId).Scan(&timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to get latest order timestamp: %w", err)
	}

	return timestamp, nil
}

// Helper function to handle null time pointers
func nullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// For better error handling, add custom error types
var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("resource not found")
)

// Add metrics tracking (optional)
type Metrics struct {
	OrdersSynced     int64
	SyncDuration     time.Duration
	BatchProcessTime time.Duration
	ErrorCount       int64
}

// Add this to your repository if you want to track metrics
func (r *PostgresRepository) GetMetrics(ctx context.Context, shopName, accountId string) (*Metrics, error) {
	query := `
		SELECT 
			COUNT(*) as orders_synced,
			MAX(updated_at) - MIN(created_at) as sync_duration
		FROM orders 
		WHERE shop_name = $1 AND account_id = $2`

	metrics := &Metrics{}
	row := r.db.QueryRowContext(ctx, query, shopName, accountId)
	if err := row.Scan(&metrics.OrdersSynced, &metrics.SyncDuration); err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return metrics, nil
}