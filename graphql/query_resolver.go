package main

import (
	"context"
	"fmt"

	"github.com/Shridhar2104/logilo/graphql/models"
)

type queryResolver struct {
    server *Server
}

func (r *queryResolver) GetAccountByID(ctx context.Context, email string, password string) (*models.Account, error) {
    accountResp, err := r.server.accountClient.LoginAndGetAccount(ctx, email, password)
    if err != nil {
        return nil, err
    }

    return &models.Account{
        ID:       accountResp.ID.String(),
        Name:     accountResp.Name,
        Password: accountResp.Password,
        Email:    accountResp.Email,
    }, nil
}

func (r *queryResolver) Accounts(ctx context.Context, pagination PaginationInput) ([]*models.Account, error) {
    res, err := r.server.accountClient.ListAccounts(ctx, uint64(pagination.Skip), uint64(pagination.Take))
    if err != nil {
        return nil, err
    }

    accounts := make([]*models.Account, len(res))
    for i, account := range res {
        accounts[i] = &models.Account{ID: account.ID.String(), Name: account.Name}
    }
    return accounts, nil
}

func (r *queryResolver) GetOrdersForAccount(ctx context.Context, accountId string, pagination *OrderPaginationInput) (*OrderConnection, error) {
    pageSize := 20
    if pagination != nil && pagination.PageSize != nil {
        pageSize = *pagination.PageSize
    }

    page := 1
    if pagination != nil && pagination.Page != nil {
        page = *pagination.Page
    }

    // Call the shopify client
    resp, err := r.server.shopifyClient.GetOrdersForAccount(ctx, accountId, int32(page), int32(pageSize))
    if err != nil {
        return nil, fmt.Errorf("failed to get orders: %w", err)
    }

    edges := make([]*OrderEdge, len(resp.Orders))
    for i, order := range resp.Orders {
        edges[i] = &OrderEdge{
            Node: &models.Order{
                ID:                fmt.Sprintf("%d", order.ID),
                Name:              order.Name,
                Amount:            order.TotalPrice,
                AccountId:         accountId,
                CreatedAt:         order.CreatedAt,
                Currency:          order.Currency,
                TotalPrice:       order.TotalPrice,
                SubtotalPrice:    order.SubtotalPrice,
                TotalTax:         &order.TotalTax,
                FinancialStatus:  order.FinancialStatus,
                FulfillmentStatus: order.FulfillmentStatus,
                Customer: &models.Customer{
                    Email:     order.Customer.Email,
                    FirstName: order.Customer.FirstName,
                    LastName:  order.Customer.LastName,
                    Phone:     order.Customer.Phone,
                },
            },
        }
    }

    return &OrderConnection{
        Edges: edges,
        PageInfo: &PageInfo{
            HasNextPage:     page < int(resp.TotalPages),
            HasPreviousPage: page > 1,
            TotalPages:      int(resp.TotalPages),
            CurrentPage:     page,
        },
        TotalCount: int(resp.TotalCount),
    }, nil
}



func (r *queryResolver) GetOrder(ctx context.Context, id string) (*models.Order, error) {
    order, err := r.server.shopifyClient.GetOrder(ctx, id)
    if err != nil {
        return nil, err
    }

    return &models.Order{
        ID:                string(order.ID),
        Name:              order.Name,
        Amount:            order.TotalPrice,
        AccountId:         "", // Need to get from context or order
        CreatedAt:         order.CreatedAt,
        Currency:          order.Currency,
        TotalPrice:        order.TotalPrice,
        SubtotalPrice:     order.SubtotalPrice,
        TotalDiscounts:    &order.TotalDiscounts,
        TotalTax:          &order.TotalTax,
        TaxesIncluded:     order.TaxesIncluded,
        FinancialStatus:   order.FinancialStatus,
        FulfillmentStatus: order.FulfillmentStatus,
        // ShopName:          order.ShopName,
        Customer: &models.Customer{
            Email:     order.Customer.Email,
            FirstName: order.Customer.FirstName,
            LastName:  order.Customer.LastName,
            Phone:     order.Customer.Phone,
        },
    }, nil
}

func (r *queryResolver) GetWalletDetails(ctx context.Context, input GetWalletDetailsInput) (*WalletDetailsResponse, error) {
    // Call the wallet client
    resp, err := r.server.paymentClient.GetWalletDetails(ctx, input.AccountID)
    if err != nil {
        return &WalletDetailsResponse{
            WalletDetails: nil,
            Errors: []*Error{{
                Code:    "WALLET_DETAILS_ERROR",
                Message: fmt.Sprintf("Failed to get wallet details: %v", err),
            }},
        }, nil
    }

    // Map the protobuf response to our GraphQL model
    return &WalletDetailsResponse{
        WalletDetails: &WalletDetails{
            AccountID:    resp.AccountId,
            Balance:      &resp.Balance,
        },
        Errors: nil,
    }, nil
}
// Ping is a simple health check method
func (r *queryResolver) Ping(ctx context.Context) (string, error) {
    return "pong", nil
}