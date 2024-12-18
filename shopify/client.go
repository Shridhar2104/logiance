package shopify

import (
	"context"
	"github.com/Shridhar2104/logilo/shopify/pb"
	"google.golang.org/grpc"
	"fmt"
)

// Client struct for gRPC communication.
type Client struct {
	conn    *grpc.ClientConn
	service pb.ShopifyServiceClient
}

// NewClient creates a new gRPC client.
func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewShopifyServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}




// GetOrdersForShopAndAccount fetches orders for a specific shop and account.
func (c *Client) GetOrdersForShopAndAccount(ctx context.Context, shopName, accountId string) ([]*Order, error) {
	res, err := c.service.GetOrdersForShopAndAccount(ctx, &pb.GetOrdersForShopAndAccountRequest{
		ShopName:  shopName,
		AccountId: accountId,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, len(res.Orders))
	for i, o := range res.Orders {
		orders[i] = &Order{
			ID:         o.Id,
			ShopName:   shopName,
			AccountId:  accountId,
			TotalPrice: float64(o.TotalPrice),
			OrderId:    o.OrderId,
		}
	}
	return orders, nil
}
//http://djcajdjd?code=jfl&shopname

// GenerateAuthURL generates an authorization URL for a Shopify store.
func (c *Client) GenerateAuthURL(ctx context.Context, shopName string) (string, error) {
	req := &pb.GetAuthorizationURLRequest{
		ShopName:  shopName,
		State: "your_unique_nonce",
		
	}
	resp, err := c.service.GetAuthorizationURL(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate auth URL: %v", err)
	}
	return resp.AuthUrl, nil
}

// ExchangeAccessToken exchanges a Shopify auth code for an access token.
func (c *Client) ExchangeAccessToken(ctx context.Context, shopName, code, accountId string) error {
	req := &pb.ExchangeAccessTokenRequest{
		ShopName:  shopName,
		Code:      code,
		AccountId: accountId,
	}
	_, err := c.service.ExchangeAccessToken(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to exchange access token: %v", err)
	}
	return nil
}

