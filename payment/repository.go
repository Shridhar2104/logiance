package payment

import (
	"context"
	"database/sql"
	"time"
)

// WalletRepository defines the interface for wallet operations.
type Repository interface {
	Close() 
	RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error)
	DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error)
	ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error)
	GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error)
}

// Implementation of WalletRepository.
type postgresRepository struct {
	db *sql.DB
}

// RemittanceDetail represents the result of processing a single order's remittance.
type RemittanceDetail struct {
	OrderID   string
	Amount    float64
	Processed bool
}

// Transaction represents a single transaction record.
type Transaction struct {
	TransactionID string
	TransactionType string
	Amount         float64
	OrderID        sql.NullString
	Timestamp      time.Time
}

// NewPostgresWalletRepository creates a new WalletRepository implementation.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// RechargeWallet adds funds to the user's wallet.
func (r *postgresRepository) RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Update wallet balance.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallets (account_id, balance, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (account_id)
		DO UPDATE SET balance = wallets.balance + $2, updated_at = NOW()`,
		accountID, amount,
	)
	if err != nil {
		return 0, err
	}

	// Fetch updated balance.
	var newBalance float64
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&newBalance)
	if err != nil {
		return 0, err
	}

	// Log the transaction.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (account_id, transaction_type, amount)
		VALUES ($1, 'recharge', $2)`,
		accountID, amount,
	)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}

// DeductBalance deducts funds from the user's wallet for shipping an order.
func (r *postgresRepository) DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Check if the wallet has sufficient balance.
	var currentBalance float64
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&currentBalance)
	if err != nil {
		return 0, err
	}
	if currentBalance < amount {
		return 0, sql.ErrNoRows // Insufficient funds.
	}

	// Deduct the balance.
	_, err = tx.ExecContext(ctx, `
		UPDATE wallets SET balance = balance - $2, updated_at = NOW()
		WHERE account_id = $1`,
		accountID, amount,
	)
	if err != nil {
		return 0, err
	}

	// Log the transaction.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (account_id, transaction_type, amount, order_id)
		VALUES ($1, 'deduction', $2, $3)`,
		accountID, amount, orderID,
	)
	if err != nil {
		return 0, err
	}

	// Fetch updated balance.
	var newBalance float64
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&newBalance)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}

// ProcessRemittance processes COD remittance for delivered orders after 15 days.
func (r *postgresRepository) ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error) {
	var details []RemittanceDetail

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, orderID := range orderIDs {
		var amount float64
		err = tx.QueryRowContext(ctx, `
			SELECT total_price FROM orders 
			WHERE id = $1 AND delivery_date <= NOW() - INTERVAL '15 days'
			AND remittance_processed = FALSE`, orderID).Scan(&amount)
		if err == sql.ErrNoRows {
			details = append(details, RemittanceDetail{OrderID: orderID, Amount: 0, Processed: false})
			continue
		} else if err != nil {
			return nil, err
		}

		// Update the wallet balance.
		_, err = tx.ExecContext(ctx, `
			UPDATE wallets SET balance = balance + $1, updated_at = NOW()
			WHERE account_id = $2`, amount, accountID)
		if err != nil {
			return nil, err
		}

		// Mark the order as remitted.
		_, err = tx.ExecContext(ctx, `
			UPDATE orders SET remittance_processed = TRUE
			WHERE id = $1`, orderID)
		if err != nil {
			return nil, err
		}

		// Log the transaction.
		_, err = tx.ExecContext(ctx, `
			INSERT INTO transactions (account_id, transaction_type, amount, order_id)
			VALUES ($1, 'remittance', $2, $3)`, accountID, amount, orderID)
		if err != nil {
			return nil, err
		}

		details = append(details, RemittanceDetail{OrderID: orderID, Amount: amount, Processed: true})
	}

	return details, nil
}

// GetWalletDetails retrieves wallet balance and transaction history for a user.
func (r *postgresRepository) GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&balance)
	if err != nil {
		return 0, nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT transaction_id, transaction_type, amount, order_id, created_at 
		FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`, accountID)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		err = rows.Scan(&txn.TransactionID, &txn.TransactionType, &txn.Amount, &txn.OrderID, &txn.Timestamp)
		if err != nil {
			return 0, nil, err
		}
		transactions = append(transactions, txn)
	}

	return balance, transactions, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}