package payment

import (
	"context"
	"google.golang.org/grpc"
	
	"github.com/Shridhar2104/logilo/payment/pb"
)

type Client struct {
	connection *grpc.ClientConn
	service    pb.PaymentServiceClient
}

// NewClient creates a new gRPC client for the Payment Service.
func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err // Return error if connection fails
	}

	client := pb.NewPaymentServiceClient(conn)
	return &Client{
		connection: conn,
		service:    client,
	}, nil
}

// Close closes the gRPC client connection.
func (c *Client) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}

// RechargeWallet calls the RechargeWallet RPC method.
func (c *Client) RechargeWallet(ctx context.Context, accountID string, amount float64) (float64, error) {
	req := &pb.RechargeWalletRequest{
		AccountId: accountID,
		Amount:    amount,
	}
	resp, err := c.service.RechargeWallet(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.NewBalance, nil
}

// DeductBalance calls the DeductBalance RPC method.
func (c *Client) DeductBalance(ctx context.Context, accountID string, amount float64, orderID string) (float64, error) {
	req := &pb.DeductBalanceRequest{
		AccountId: accountID,
		Amount:    amount,
		OrderId:   orderID,
	}
	resp, err := c.service.DeductBalance(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.NewBalance, nil
}

// ProcessRemittance calls the ProcessRemittance RPC method.
func (c *Client) ProcessRemittance(ctx context.Context, accountID string, orderIDs []string) ([]*pb.RemittanceDetail, error) {
	req := &pb.ProcessRemittanceRequest{
		AccountId: accountID,
		OrderIds:  orderIDs,
	}
	resp, err := c.service.ProcessRemittance(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.RemittanceDetails, nil
}

// GetWalletDetails calls the GetWalletDetails RPC method.
func (c *Client) GetWalletDetails(ctx context.Context, accountID string) (*pb.WalletDetailsResponse, error) {
	req := &pb.GetWalletDetailsRequest{
		AccountId: accountID,
	}
	return c.service.GetWalletDetails(ctx, req)
}
