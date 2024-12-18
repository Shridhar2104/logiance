package shopify

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ShopifyRepo combines the repository for tokens and orders.
type Repository interface {
	Close()
	SaveShopCredentials(ctx context.Context, shopName, accountId, token string) error
	GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error)
}

// postgresRepository implements ShopifyRepo using PostgreSQL.
type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository initializes a new PostgreSQL repository.
func NewPostgresRepository(url string) (*postgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

// SaveShopCredentials stores or updates shop credentials (token).
func (r *postgresRepository) SaveShopCredentials(ctx context.Context, shopName, accountId, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tokens (shop_name, account_id, token) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (shop_name, account_id) 
		DO UPDATE SET token = $3`,
		shopName, accountId, token,
	)
	if err != nil {
		return fmt.Errorf("failed to save shop credentials: %w", err)
	}
	return nil
}

func (r *postgresRepository) Close() {
	r.db.Close()

}


func (r *postgresRepository) GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, updated_at, shop_name, account_id, order_id, total_price 
		FROM orders WHERE shop_name = $1 AND account_id = $2`,
		shopName, accountId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.ShopName, &order.AccountId, &order.OrderId, &order.TotalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

