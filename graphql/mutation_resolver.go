package main

import (
	"context"
    "fmt"
    "github.com/google/uuid"
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
// Shipping mutations
func (r *mutationResolver) CalculateShippingRates(ctx context.Context, input ShippingRateInput) (*models.ShippingRateResponse, error) {
    return r.server.Shipping().CalculateShippingRates(ctx, input)
}

func (r *mutationResolver) GetAvailableCouriers(ctx context.Context, input AvailabilityInput) (*models.CourierAvailabilityResponse, error) {
    return r.server.Shipping().GetAvailableCouriers(ctx, input)
}

func (r *mutationResolver) CreateShipment(ctx context.Context, input CreateShipmentInput) (*ShipmentResponse, error) {
    return r.server.Shipping().CreateShipment(ctx, input)
}

func (r *mutationResolver) TrackShipment(ctx context.Context, input TrackingInput) (*TrackingResponse, error) {
    return r.server.Shipping().TrackShipment(ctx, input)
}



// Add new bank account mutations
func (r *mutationResolver) AddBankAccount(ctx context.Context, userID string, input BankAccountInput) (*BankAccount, error) {
    bankAccount := &account.BankAccount{
        UserID:          userID,
        AccountNumber:   input.AccountNumber,
        AccountType:     input.AccountType,     // New field
        BranchName:      input.BranchName,      // New field
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
        AccountType:     resp.AccountType,     // New field
        BranchName:      resp.BranchName,      // New field
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
        AccountType:     input.AccountType,     // New field
        BranchName:      input.BranchName,      // New field
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
        AccountType:     resp.AccountType,     // New field
        BranchName:      resp.BranchName,      // New field
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

 func (r *mutationResolver) AddWareHouse(ctx context.Context, userID string, input WareHouseInput) (*WareHouse, error) {
    // Handle nil Landmark properly
    var landmark string
    if input.Landmark != nil {
        landmark = *input.Landmark
    }

    wh := &account.Address{
        UserID:          userID,
        ContactPerson:   input.ContactPerson,
        ContactNumber:   input.ContactNumber,
        EmailAddress:    input.EmailAddress,
        CompleteAddress: input.CompleteAddress,
        Landmark:        landmark,
        Pincode:        input.Pincode,
        City:           input.City,
        State:          input.State,
        Country:        input.Country,
    }
    
    resp, err := r.server.accountClient.AddAddress(ctx, wh)
    if err != nil {
        return nil, fmt.Errorf("failed to add warehouse: %w", err)
    }

    landmark = resp.Landmark
    return &WareHouse{
        ID:              resp.ID.String(),
        UserID:          resp.UserID,
        ContactPerson:   resp.ContactPerson,
        ContactNumber:   resp.ContactNumber,
        EmailAddress:    resp.EmailAddress,
        CompleteAddress: resp.CompleteAddress,
        Landmark:        &landmark,
        Pincode:        resp.Pincode,
        City:           resp.City,
        State:          resp.State,
        Country:        resp.Country,
        // CreatedAt:      resp.CreatedAt,
        // UpdatedAt:      resp.UpdatedAt,
    }, nil
}

func (r *mutationResolver) UpdateWareHouse(ctx context.Context, id string, input WareHouseInput) (*WareHouse, error) {
    // Handle nil Landmark properly
    var landmark string
    if input.Landmark != nil {
        landmark = *input.Landmark
    }

    wh := &account.Address{
        ID:              uuid.MustParse(id),
        ContactPerson:   input.ContactPerson,
        ContactNumber:   input.ContactNumber,
        EmailAddress:    input.EmailAddress,
        CompleteAddress: input.CompleteAddress,
        Landmark:        landmark,
        Pincode:        input.Pincode,
        City:           input.City,
        State:          input.State,
        Country:        input.Country,
    }

    resp, err := r.server.accountClient.UpdateAddress(ctx, wh)
    if err != nil {
        return nil, fmt.Errorf("failed to update warehouse: %w", err)
    }

    landmark = resp.Landmark
    return &WareHouse{
        ID:              resp.ID.String(),
        UserID:          resp.UserID,
        ContactPerson:   resp.ContactPerson,
        ContactNumber:   resp.ContactNumber,
        EmailAddress:    resp.EmailAddress,
        CompleteAddress: resp.CompleteAddress,
        Landmark:        &landmark,
        Pincode:        resp.Pincode,
        City:           resp.City,
        State:          resp.State,
        Country:        resp.Country,
        // CreatedAt:      resp.CreatedAt,
        // UpdatedAt:      resp.UpdatedAt,
    }, nil
}

func (r *mutationResolver) DeleteWareHouse(ctx context.Context, id string) (bool, error) {
    err := r.server.accountClient.DeleteAddress(ctx, id)
    if err != nil {
        return false, fmt.Errorf("failed to delete warehouse: %w", err)
    }
    return true, nil
}