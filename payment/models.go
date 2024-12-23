package payment

import (
	"database/sql"
	"time"
)

type Transaction struct {
	TransactionID   string        `json:"transaction_id"`
	TransactionType string        `json:"transaction_type"`
	Amount          float64       `json:"amount"`
	OrderID         sql.NullString `json:"order_id"`
	Timestamp       time.Time     `json:"timestamp"`
}

type WalletDetails struct {
	AccountID    string        `json:"account_id"`
	Balance      float64       `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

type RemittanceDetail struct {
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	Processed bool    `json:"processed"`
}
