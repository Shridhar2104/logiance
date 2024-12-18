package payment

import (
	"context"
	"fmt"
	"net"

	"github.com/Shridhar2104/logilo/payment/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service Service
}

// NewGRPCServer creates and starts a new gRPC server.
func NewGRPCServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterPaymentServiceServer(server, &grpcServer{
		service: service,
	})
	reflection.Register(server)
	fmt.Printf("gRPC server running on port %d\n", port)
	return server.Serve(lis)
}

// RechargeWallet handles wallet recharge requests.
func (s *grpcServer) RechargeWallet(ctx context.Context, req *pb.RechargeRequest) (*pb.RechargeResponse, error) {
	newBalance, err := s.service.RechargeWallet(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &pb.RechargeResponse{
		NewBalance: newBalance,
	}, nil
}

// DeductBalance handles wallet balance deduction for an order.
func (s *grpcServer) DeductBalance(ctx context.Context, req *pb.DeductionRequest) (*pb.DeductionResponse, error) {
	newBalance, err := s.service.DeductBalance(ctx, req.UserId, req.Amount, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &pb.DeductionResponse{
		NewBalance: newBalance,
	}, nil
}

// ProcessRemittance handles the remittance processing for COD orders.
func (s *grpcServer) ProcessRemittance(ctx context.Context, req *pb.RemittanceRequest) (*pb.RemittanceResponse, error) {
	remittanceDetails, err := s.service.ProcessRemittance(ctx, req.UserId, req.OrderIds)
	if err != nil {
		return nil, err
	}

	// Map the remittance details to the gRPC response format.
	var remittanceItems []*pb.RemittanceDetail
	for _, detail := range remittanceDetails {
		remittanceItems = append(remittanceItems, &pb.RemittanceDetail{
			OrderId:         detail.OrderID,
			Amount:          detail.Amount,
			Processed:       detail.Processed,
		})
	}

	return &pb.RemittanceResponse{
		Details: remittanceItems,
	}, nil
}

// GetWalletDetails retrieves wallet balance and transaction history.
func (s *grpcServer) GetWalletDetails(ctx context.Context, req *pb.WalletDetailsRequest) (*pb.WalletDetailsResponse, error) {
	balance, transactions, err := s.service.GetWalletDetails(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Map the transactions to the gRPC response format.
	var transactionItems []*pb.Transaction
	for _, t := range transactions {
		transactionItems = append(transactionItems, &pb.Transaction{
			TransactionId:   t.TransactionID,
			TransactionType: t.TransactionType,
			Amount:          t.Amount,
			OrderId:         t.OrderID.String,
		})
	}

	return &pb.WalletDetailsResponse{
		Balance:      balance,
		TransactionHistory: transactionItems,

	}, nil
}
