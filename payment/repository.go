package payment

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error)
	DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error)
	ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error)
	GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

//recharge wallet 
func (r *repository) RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var balance float64
	err = tx.QueryRowContext(ctx, `UPDATE wallets SET balance = balance + $1 WHERE account_id = $2 RETURNING balance`, amount, accountID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO transactions (transaction_type, amount, account_id, timestamp) VALUES ('recharge', $1, $2, NOW())`, amount, accountID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return balance, tx.Commit()
}


//deduct balance
func (r *repository) DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var balance float64
	err = tx.QueryRowContext(ctx, `SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if balance < amount {
		tx.Rollback()
		return 0, errors.New("insufficient balance")
	}

	err = tx.QueryRowContext(ctx, `UPDATE wallets SET balance = balance - $1 WHERE account_id = $2 RETURNING balance`, amount, accountID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO transactions (transaction_type, amount, order_id, account_id, timestamp) VALUES ('deduction', $1, $2, $3, NOW())`, amount, orderID, accountID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return balance, tx.Commit()
}

func (r *repository) ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error) {
	var details []RemittanceDetail
	for _, orderID := range orderIDs {
		var amount float64
		err := r.db.QueryRowContext(ctx, `SELECT amount FROM orders WHERE order_id = $1 AND account_id = $2`, orderID, accountID).Scan(&amount)
		if err != nil {
			details = append(details, RemittanceDetail{OrderID: orderID, Processed: false})
			continue
		}

		_, err = r.DeductBalance(ctx, accountID, amount, orderID)
		if err != nil {
			details = append(details, RemittanceDetail{OrderID: orderID, Processed: false})
		} else {
			details = append(details, RemittanceDetail{OrderID: orderID, Amount: amount, Processed: true})
		}
	}
	return details, nil
}

//wallet details
func (r *repository) GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx, `SELECT balance FROM wallets WHERE account_id = $1`, accountID).Scan(&balance)
	if err != nil {
		return 0, nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT transaction_id, transaction_type, amount, order_id, timestamp FROM transactions WHERE account_id = $1`, accountID)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err = rows.Scan(&t.TransactionID, &t.TransactionType, &t.Amount, &t.OrderID, &t.Timestamp)
		if err != nil {
			return 0, nil, err
		}
		transactions = append(transactions, t)
	}

	return balance, transactions, nil
}
