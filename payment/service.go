package payment

import (
    "context"
    "errors"
    "github.com/razorpay/razorpay-go"
)

type Service interface {
    RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error)
    DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error)
    ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error)
    GetWalletDetails(ctx context.Context, accountID string) (WalletDetails, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error) {
    if amount <= 0 {
        return 0, errors.New("amount must be greater than zero")
    }

    razorpayClient := razorpay.NewClient("YOUR_RAZORPAY_API_KEY", "YOUR_RAZORPAY_SECRET")

    // Create order data
    orderData := map[string]interface{}{
        "amount":   int(amount * 100), // Amount in paise
        "currency": "INR",
        "receipt":  "txn_" + accountID,
    }

    // Create order options
	// Abhi k liye empty 
    orderOptions := map[string]string{}

    // Create an order in Razorpay
    order, err := razorpayClient.Order.Create(orderData, orderOptions)
    if err != nil {
        return 0, errors.New("failed to create Razorpay order: " + err.Error())
    }

    // Check if the order was created successfully
    if order["id"] == nil {
        return 0, errors.New("order creation failed")
    }

    // Simulate payment confirmation (this part is typically handled by Razorpay's frontend)
    // Replace with actual payment confirmation logic
    paymentConfirmed := true // Assume payment is confirmed for now
    if !paymentConfirmed {
        return 0, errors.New("payment not confirmed")
    }

    // Update wallet balance in the repository
    return s.repo.RechargeWallet(ctx, accountID, amount)
}

func (s *service) DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error) {
    if amount <= 0 {
        return 0, errors.New("amount must be greater than zero")
    }
    if orderID == "" {
        return 0, errors.New("orderID is required")
    }
    return s.repo.DeductBalance(ctx, accountID, amount, orderID)
}

func (s *service) ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error) {
    if len(orderIDs) == 0 {
        return nil, errors.New("orderIDs cannot be empty")
    }
    return s.repo.ProcessRemittance(ctx, accountID, orderIDs)
}

func (s *service) GetWalletDetails(ctx context.Context, accountID string) (WalletDetails, error) {
    balance, transactions, err := s.repo.GetWalletDetails(ctx, accountID)
    if err != nil {
        return WalletDetails{}, err
    }
    return WalletDetails{
        AccountID:    accountID,
        Balance:      balance,
        Transactions: transactions,
    }, nil
}