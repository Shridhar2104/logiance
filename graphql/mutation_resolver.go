package main

import (
	"context"
    "fmt"

	"github.com/Shridhar2104/logilo/account"
	"github.com/Shridhar2104/logilo/graphql/models"
)

type mutationResolver struct {
	server *Server
}


func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*models.Account, error) {

	a:= &account.Account{
		Name: input.Name,
		Password: input.Password,
		Email: input.Email,
	}

	res, err := r.server.accountClient.CreateAccount(ctx, a)
	if err != nil {
		return nil, err
	}

	return &models.Account{
		ID: res.ID.String(),
		Name: res.Name,
		Password: res.Password,
		Email: res.Email,
		Orders: nil,
		ShopNames: nil,
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

func (r *mutationResolver) SyncOrders(ctx context.Context, accountId string) (*models.SyncOrdersResult, error) {
    shopResults, err := r.server.shopifyClient.SyncOrders(ctx, accountId)
    if err != nil {
        errorMsg := fmt.Sprintf("Failed to sync orders: %v", err)
        return &models.SyncOrdersResult{
            OverallSuccess: false,
            Message:       errorMsg,
            ShopResults:   nil,
        }, nil
    }

    shopDetailsSlice := make([]*models.ShopSyncDetails, 0, len(shopResults))
    allSuccessful := true

    for shopName, status := range shopResults {
        if !status.Success {
            allSuccessful = false
        }

        shopDetailsSlice = append(shopDetailsSlice, &models.ShopSyncDetails{
            ShopName: shopName,
            Status: &models.ShopSyncStatus{
                Success:      status.Success,
                ErrorMessage: status.ErrorMessage,
                OrdersCount: int(status.OrdersSynced),
            },
        })
    }

    return &models.SyncOrdersResult{
        OverallSuccess: allSuccessful,
        Message:       "Order synchronization completed",
        ShopResults:   shopDetailsSlice,
    }, nil
}






//recharge wallet
func (r *mutationResolver) RechargeWallet(ctx context.Context, input RechargeWalletInput) (*WalletOperationResponse, error) {
    newBalance, err := r.server.paymentClient.RechargeWallet(ctx, input.AccountID, input.Amount)
    if err != nil {
        return &WalletOperationResponse{
            NewBalance: 0,
            Errors: []*Error{{
                Code:    "RECHARGE_FAILED",
                Message: fmt.Sprintf("Failed to recharge wallet: %v", err),
            }},
        }, nil
    }

    return &WalletOperationResponse{
        NewBalance: newBalance,
        Errors:    nil,
    }, nil
}

//deduct balance
func (r *mutationResolver) DeductBalance(ctx context.Context, input DeductBalanceInput) (*WalletOperationResponse, error) {
    newBalance, err := r.server.paymentClient.DeductBalance(ctx, input.AccountID, input.Amount, input.OrderID)
    if err != nil {
        return &WalletOperationResponse{
            NewBalance: 0,
            Errors: []*Error{{
                Code:    "DEDUCTION_FAILED",
                Message: fmt.Sprintf("Failed to deduct balance: %v", err),
            }},
        }, nil
    }

    return &WalletOperationResponse{
        NewBalance: newBalance,
        Errors:    nil,
    }, nil
}