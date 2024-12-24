package models

import "time"

// Add BankAccount model
type BankAccount struct {
    UserID          string    `json:"userId"`
    AccountNumber   string    `json:"accountNumber"`
    BeneficiaryName string    `json:"beneficiaryName"`
    IFSCCode        string    `json:"ifscCode"`
    BankName        string    `json:"bankName"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}

// Add BankAccountInput model
type BankAccountInput struct {
    AccountNumber   string `json:"accountNumber"`
    BeneficiaryName string `json:"beneficiaryName"`
    IFSCCode        string `json:"ifscCode"`
    BankName        string `json:"bankName"`
}

// Update Account model to include BankAccount
type Account struct {
    ID         string      `json:"id"`
    Name       string      `json:"name"`
    Password   string      `json:"password"`
    Email      string      `json:"email"`
    Orders     []Order     `json:"orders"`
    ShopNames  []ShopName  `json:"shopnames"`
    BankAccount *BankAccount `json:"bankAccount,omitempty"`
}


type ShopName struct {
    Shopname string `json:"shopname"`
}

type Order struct {
    ID                string     `json:"id"`
    Name             string     `json:"name"`
    OrderNumber      *int       `json:"orderNumber"`
    Amount           float64    `json:"amount"`
    AccountId        string     `json:"accountId"`
    CreatedAt        time.Time  `json:"createdAt"`
    UpdatedAt        time.Time  `json:"updatedAt"`
    CancelledAt      *time.Time `json:"cancelledAt,omitempty"`
    ClosedAt         *time.Time `json:"closedAt,omitempty"`
    ProcessedAt      *time.Time `json:"processedAt,omitempty"`
    Currency         string     `json:"currency"`
    TotalPrice       float64    `json:"totalPrice"`
    SubtotalPrice    float64    `json:"subtotalPrice"`
    TotalDiscounts   *float64   `json:"totalDiscounts,omitempty"`
    TotalTax         *float64   `json:"totalTax,omitempty"`
    TaxesIncluded    bool       `json:"taxesIncluded"`
    FinancialStatus  string     `json:"financialStatus"`
    FulfillmentStatus string    `json:"fulfillmentStatus"`
    Description      string     `json:"description"`
    LineItems        []*OrderLineItem `json:"lineItems"`
    Customer         *Customer  `json:"customer,omitempty"`
    ShopName         string     `json:"shopName"`
}

type OrderLineItem struct {
    ID          string  `json:"id"`
    Amount      float64 `json:"amount"`
    Description string  `json:"description"`
}

type Customer struct {
    ID        string `json:"id,omitempty"`
    Email     string `json:"email,omitempty"`
    FirstName string `json:"firstName,omitempty"`
    LastName  string `json:"lastName,omitempty"`
    Phone     string `json:"phone,omitempty"`
}

type OrderConnection struct {
    Edges      []*OrderEdge `json:"edges"`
    PageInfo   *PageInfo    `json:"pageInfo"`
    TotalCount int          `json:"totalCount"`
}

type OrderEdge struct {
    Node   *Order `json:"node"`
    Cursor string `json:"cursor"`
}

type PageInfo struct {
    HasNextPage     bool    `json:"hasNextPage"`
    HasPreviousPage bool    `json:"hasPreviousPage"`
    StartCursor     *string `json:"startCursor,omitempty"`
    EndCursor       *string `json:"endCursor,omitempty"`
    TotalPages      int     `json:"totalPages"`
    CurrentPage     int     `json:"currentPage"`
}

type OrderPaginationInput struct {
    Page     *int         `json:"page,omitempty"`
    PageSize *int         `json:"pageSize,omitempty"`
    Filter   *OrderFilter `json:"filter,omitempty"`
    Sort     *OrderSort   `json:"sort,omitempty"`
}

type OrderFilter struct {
    CreatedAtStart    *time.Time `json:"createdAtStart,omitempty"`
    CreatedAtEnd      *time.Time `json:"createdAtEnd,omitempty"`
    FinancialStatus   *string    `json:"financialStatus,omitempty"`
    FulfillmentStatus *string    `json:"fulfillmentStatus,omitempty"`
    MinTotalPrice     *float64   `json:"minTotalPrice,omitempty"`
    MaxTotalPrice     *float64   `json:"maxTotalPrice,omitempty"`
    SearchTerm        *string    `json:"searchTerm,omitempty"`
}

type OrderSort struct {
    Field     OrderSortField `json:"field"`
    Direction SortDirection  `json:"direction"`
}

type OrderSortField string
const (
    OrderSortFieldCreatedAt   OrderSortField = "CREATED_AT"
    OrderSortFieldUpdatedAt   OrderSortField = "UPDATED_AT"
    OrderSortFieldOrderNumber OrderSortField = "ORDER_NUMBER"
    OrderSortFieldTotalPrice  OrderSortField = "TOTAL_PRICE"
)

type SortDirection string
const (
    SortDirectionAsc  SortDirection = "ASC"
    SortDirectionDesc SortDirection = "DESC"
)

type AccountInput struct {
    Name     string `json:"name"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

type OrderInput struct {
    AccountId string             `json:"accountId"`
    LineItems []*OrderLineItemInput `json:"lineItems"`
}

type OrderLineItemInput struct {
    ID          string  `json:"id"`
    Amount      float64 `json:"amount"`
    Description string  `json:"description"`
}

type ShopSyncStatus struct {
    Success      bool   `json:"success"`
    ErrorMessage string `json:"errorMessage,omitempty"`
    OrdersCount  int    `json:"ordersCount"`
}

type ShopSyncDetails struct {
    ShopName string         `json:"shopName"`
    Status   *ShopSyncStatus `json:"status"`
}

type SyncOrdersResult struct {
    OverallSuccess bool              `json:"overallSuccess"`
    Message        string            `json:"message,omitempty"`
    ShopResults    []*ShopSyncDetails `json:"shopResults"`
}