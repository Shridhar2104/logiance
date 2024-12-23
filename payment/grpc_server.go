package payment

import (
	"context"
	"github.com/Shridhar2104/logilo/payment/pb"
	"time"
)

type grpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service Service
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

func NewGRPCServer(service Service) pb.PaymentServiceServer {
	return &grpcServer{service: service}
}

func (s *grpcServer) RechargeWallet(ctx context.Context, req *pb.RechargeWalletRequest) (*pb.RechargeWalletResponse, error) {
	balance, err := s.service.RechargeWallet(ctx, req.AccountId, req.Amount)
	if err != nil {
		return nil, err
	}
	return &pb.RechargeWalletResponse{
		AccountId:  req.AccountId,
		NewBalance: balance,
	}, nil
}

func (s *grpcServer) DeductBalance(ctx context.Context, req *pb.DeductBalanceRequest) (*pb.DeductBalanceResponse, error) {
	balance, err := s.service.DeductBalance(ctx, req.AccountId, req.Amount, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.DeductBalanceResponse{
		AccountId:  req.AccountId,
		NewBalance: balance,
	}, nil
}

func (s *grpcServer) ProcessRemittance(ctx context.Context, req *pb.ProcessRemittanceRequest) (*pb.ProcessRemittanceResponse, error) {
	details, err := s.service.ProcessRemittance(ctx, req.AccountId, req.OrderIds)
	if err != nil {
		return nil, err
	}

	var pbDetails []*pb.RemittanceDetail
	for _, detail := range details {
		pbDetails = append(pbDetails, &pb.RemittanceDetail{
			OrderId:   detail.OrderID,
			Amount:    detail.Amount,
			Processed: detail.Processed,
		})
	}

	return &pb.ProcessRemittanceResponse{
		RemittanceDetails: pbDetails,
	}, nil
}

func (s *grpcServer) GetWalletDetails(ctx context.Context, req *pb.GetWalletDetailsRequest) (*pb.WalletDetailsResponse, error) {
	details, err := s.service.GetWalletDetails(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	var pbTransactions []*pb.Transaction
	for _, t := range details.Transactions {
		var orderID *string
		if t.OrderID.Valid {
			orderID = stringPtr(t.OrderID.String) // Convert to a pointer
		}

		// Convert time.Time to string using RFC3339 format
		timestampStr := t.Timestamp.Format(time.RFC3339)

		pbTransactions = append(pbTransactions, &pb.Transaction{
			TransactionId:   t.TransactionID,
			TransactionType: t.TransactionType,
			Amount:          t.Amount,
			OrderId:         orderID,      // Assign *string
			Timestamp:       timestampStr, // Assign formatted string
		})
	}

	return &pb.WalletDetailsResponse{
		AccountId:    details.AccountID,
		Balance:      details.Balance,
		Transactions: pbTransactions,
	}, nil
}
