package main

import (
	"context"
	"time"

	"github.com/Shridhar2104/logilo/account"
	"github.com/Shridhar2104/logilo/graphql/models"
)

type mutationResolver struct {
	server *Server
}


func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*models.Account, error) {

	a:= &account.Account{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		
	}

	res, err := r.server.accountClient.CreateAccount(ctx, a)
	if err != nil {
		return nil, err
	}

	return &models.Account{
		ID: res.ID.String(),
		Name: res.Name,
		Email: res.Email,
		Password: res.Password,
	}, nil
}


// mutation_resolver.go

func (r *mutationResolver) IntegrateShop(ctx context.Context, shopName string) (string, error) {
    // Call the Shopify client to get the authorization URL
    url, err := r.server.shopifyClient.GenerateAuthURL(ctx, shopName)
    if err != nil {
        return "", err
    }
    return url, nil
}

func (r *mutationResolver) ExchangeAccessToken(ctx context.Context, shopName, code, accountId string) (bool, error) {
	err := r.server.shopifyClient.ExchangeAccessToken(ctx, shopName, code, accountId)
	if err != nil {
		return false, err
	}
	return true, nil
}
