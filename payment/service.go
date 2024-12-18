package payment

import "context"
type Service interface {
	RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error)
	DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error)
	ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error)
	GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error)
}


type paymentService struct {
	repo Repository
}

func NewPaymentService(repo Repository) Service {
	return &paymentService{repo}
}

func (s *paymentService) RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error) {
	return s.repo.RechargeWallet(ctx, accountID, amount)
}

func (s *paymentService) DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error) {
	return s.repo.DeductBalance(ctx, accountID, amount, orderID)
}

func (s *paymentService) ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]RemittanceDetail, error) {
	return s.repo.ProcessRemittance(ctx, accountID, orderIDs)
}

func (s *paymentService) GetWalletDetails(ctx context.Context, accountID string) (float64, []Transaction, error) {
	return s.repo.GetWalletDetails(ctx, accountID)
}
