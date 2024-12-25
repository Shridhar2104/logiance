package main

import (
    "context"
    "github.com/Shridhar2104/logilo/graphql/models"
)

type accountResolver struct {
    server *Server
}


// Add this function to your existing accountResolver
func (r *accountResolver) BankAccount(ctx context.Context, obj *models.Account) (*BankAccount, error) {
    if obj.ID == "" {
        return nil, nil
    }
    
    bankAccount, err := r.server.accountClient.GetBankAccount(ctx, obj.ID)
    if err != nil {
        return nil, err
    }
    
    return &BankAccount{
        UserID:          bankAccount.UserID,
        AccountNumber:   bankAccount.AccountNumber,
        BeneficiaryName: bankAccount.BeneficiaryName,
        IfscCode:        bankAccount.IFSCCode,
        BankName:        bankAccount.BankName,
        // CreatedAt:       bankAccount.CreatedAt,
        // UpdatedAt:       bankAccount.UpdatedAt,
    }, nil
}

// Orders fetch orders for all shopnames for the given account with pagination
func (r *accountResolver) Orders(ctx context.Context, obj *models.Account) ([]*models.Order, error) {
    // Use the shopify client to get orders for this account
    resp, err := r.server.shopifyClient.GetOrdersForAccount(
        ctx,
        obj.ID,
        1,  // default page
        100, // default page size
    )
    if err != nil {
        return nil, err
    }

    orders := make([]*models.Order, len(resp.Orders))
    for i, order := range resp.Orders {
        orders[i] = &models.Order{
            ID:                string(order.ID),
            Name:              order.Name,
            Amount:            order.TotalPrice,
            AccountId:         obj.ID,
            CreatedAt:         order.CreatedAt,
            Currency:          order.Currency,
            TotalPrice:        order.TotalPrice,
            SubtotalPrice:     order.SubtotalPrice,
            // TotalDiscounts:    order.TotalDiscounts,
            // TotalTax:          order.TotalTax,
            TaxesIncluded:     order.TaxesIncluded,
            FinancialStatus:   order.FinancialStatus,
            FulfillmentStatus: order.FulfillmentStatus,
            //ShopName:          order.ShopName,
            Customer: &models.Customer{
                Email:     order.Customer.Email,
                FirstName: order.Customer.FirstName,
                LastName:  order.Customer.LastName,
                Phone:     order.Customer.Phone,
            },
        }
    }

    return orders, nil
}

// Shopnames returns the shopnames related to the account
func (r *accountResolver) Shopnames(ctx context.Context, obj *models.Account) ([]*models.ShopName, error) {
    var shopNames []*models.ShopName
    for _, shopName := range obj.ShopNames {
        shopNames = append(shopNames, &models.ShopName{
            Shopname: shopName.Shopname,
        })
    }
    return shopNames, nil
}