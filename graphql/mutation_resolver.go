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

// Add new bank account mutations
func (r *mutationResolver) AddBankAccount(ctx context.Context, userID string, input BankAccountInput) (*BankAccount, error) {
    bankAccount := &account.BankAccount{
        UserID:          userID,
        AccountNumber:   input.AccountNumber,
        BeneficiaryName: input.BeneficiaryName,
        IFSCCode:        input.IfscCode,
        BankName:        input.BankName,
    }
 
    resp, err := r.server.accountClient.AddBankAccount(ctx, bankAccount)
    if err != nil {
        return nil, fmt.Errorf("failed to add bank account: %w", err)
    }
 
    return &BankAccount{
        UserID:          resp.UserID,
        AccountNumber:   resp.AccountNumber,
        BeneficiaryName: resp.BeneficiaryName,
        IfscCode:        resp.IFSCCode,
        BankName:        resp.BankName,
        // CreatedAt:       resp.CreatedAt,
        // UpdatedAt:       resp.UpdatedAt,
    }, nil
 }
 
 func (r *mutationResolver) UpdateBankAccount(ctx context.Context, userID string, input BankAccountInput) (*BankAccount, error) {
    bankAccount := &account.BankAccount{
        UserID:          userID,
        AccountNumber:   input.AccountNumber,
        BeneficiaryName: input.BeneficiaryName,
        IFSCCode:        input.IfscCode,
        BankName:        input.BankName,
    }
 
    resp, err := r.server.accountClient.UpdateBankAccount(ctx, bankAccount)
    if err != nil {
        return nil, fmt.Errorf("failed to update bank account: %w", err)
    }
 
    return &BankAccount{
        UserID:          resp.UserID,
        AccountNumber:   resp.AccountNumber,
        BeneficiaryName: resp.BeneficiaryName,
        IfscCode:        resp.IFSCCode,
        BankName:        resp.BankName,
        // CreatedAt:       resp.CreatedAt,
        // UpdatedAt:       resp.UpdatedAt,
    }, nil
 }
 
 func (r *mutationResolver) DeleteBankAccount(ctx context.Context, userID string) (bool, error) {
    err := r.server.accountClient.DeleteBankAccount(ctx, userID)
    if err != nil {
        return false, fmt.Errorf("failed to delete bank account: %w", err)
    }
    return true, nil
 }