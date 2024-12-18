package payment

import (
	"context"

	"google.golang.org/grpc"
	"github.com/Shridhar2104/logilo/payment/pb"
)


type Client struct {
	conn *grpc.ClientConn
	service pb.PaymentServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c:= pb.NewPaymentServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) RechargeWallet(ctx context.Context, userId string, amount float64) (float64, error) {
	
}

func (c *Client) DeductBalance(ctx context.Context, userId string, amount float64, orderID string) (float64, error) {

}
func (c *Client) ProcessRemittance(ctx context.Context, userId string, amount float64, orderID string) (float64, error) {

}


func (c *Client) GetWalletDetails(ctx context.Context, userId string) (*pb.WalletDetailsResponse, error) {

}