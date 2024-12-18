package shopify

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)
type Service interface{
	GenerateAuthURL(ctx context.Context, shopName, state string)(string, error)
	ExchangeAccessToken(ctx context.Context, shop, code , accountId string) error
	GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error)
}

type ShopifyService struct {
	App goshopify.App
	Repo Repository // Interface for DB operations
}

func (s *ShopifyService) GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error) {

	return s.Repo.GetOrdersForShopAndAccount(ctx, shopName, accountId)

}

func NewShopifyService(apiKey, apiSecret, redirectURL string, repo Repository) Service {
	return &ShopifyService{
		App: goshopify.App{
			ApiKey:      apiKey,
			ApiSecret:   apiSecret,
			RedirectUrl: redirectURL,
			Scope: "read_products,read_orders,read_fulfillments,read_all_orders,read_merchant_managed_fulfillment_orders,write_merchant_managed_fulfillment_orders,write_fulfillments",
		},
		Repo: repo,
	}
}



// GenerateAuthURL generates the Shopify authorization URL
func (s *ShopifyService) GenerateAuthURL(ctx context.Context, shopName, state string) (string, error) {
	return s.App.AuthorizeUrl(shopName, state)
}

// ExchangeAccessToken handles the callback and exchanges code for access token
func (s *ShopifyService) ExchangeAccessToken(ctx context.Context, shop, code , accountId string) error {
	accessToken, err := s.App.GetAccessToken(ctx, shop, code)
	if err != nil {
		return err
	}
	// Save to database
	return s.Repo.SaveShopCredentials(ctx, shop, accessToken, accountId)
}
