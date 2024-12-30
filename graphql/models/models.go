package models

import "time"

// Add BankAccount model
type BankAccount struct {
	UserID          string `json:"userId"`
	AccountNumber   string `json:"accountNumber"`
	BeneficiaryName string `json:"beneficiaryName"`
	IfscCode        string `json:"ifscCode"`
	BankName        string `json:"bankName"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}


// Add BankAccountInput model
type BankAccountInput struct {
	AccountNumber   string `json:"accountNumber"`
	BeneficiaryName string `json:"beneficiaryName"`
	IfscCode        string `json:"ifscCode"`
	BankName        string `json:"bankName"`
}
// WareHouse
type WareHouse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userId"`
	ContactPerson   string  `json:"contactPerson"`
	ContactNumber   string  `json:"contactNumber"`
	EmailAddress    string  `json:"emailAddress"`
	CompleteAddress string  `json:"completeAddress"`
	Landmark        *string `json:"landmark,omitempty"`
	Pincode         string  `json:"pincode"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	Country         string  `json:"country"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type WareHouseInput struct {
	ContactPerson   string  `json:"contactPerson"`
	ContactNumber   string  `json:"contactNumber"`
	EmailAddress    string  `json:"emailAddress"`
	CompleteAddress string  `json:"completeAddress"`
	Landmark        *string `json:"landmark,omitempty"`
	Pincode         string  `json:"pincode"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	Country         string  `json:"country"`
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
    WareHouses  []*WareHouse  `json:"warehouses"`
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



type WalletStatus string

const (
    WalletStatusActive    WalletStatus = "ACTIVE"
    WalletStatusInactive  WalletStatus = "INACTIVE"
    WalletStatusSuspended WalletStatus = "SUSPENDED"
    WalletStatusClosed    WalletStatus = "CLOSED"
)

type WalletDetails struct {
    AccountID    string       `json:"accountId"`
    Balance      float64      `json:"balance"`
    Currency     string       `json:"currency"`
    Status       WalletStatus `json:"status"`
    LastUpdated  time.Time    `json:"lastUpdated"`
}

type GetWalletDetailsInput struct {
    AccountID string `json:"accountId"`
}

type WalletDetailsResponse struct {
    WalletDetails *WalletDetails `json:"walletDetails"`
    Errors        []*Error       `json:"errors"`
}

type Error struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}


type RechargeWalletInput struct {
    AccountID string  `json:"accountId"`
    Amount    float64 `json:"amount"`
}

type DeductBalanceInput struct {
    AccountID string  `json:"accountId"`
    Amount    float64 `json:"amount"`
    OrderID   string  `json:"orderId"`
}


type WalletOperationResponse struct {
    NewBalance float64  `json:"newBalance"`
    Errors     []*Error `json:"errors"`
}


// Enums
type PaymentMode string
const (
    PaymentModeCOD     PaymentMode = "COD"
    PaymentModePrepaid PaymentMode = "PREPAID"
)

type PaymentType string
const (
    PaymentTypeCOD     PaymentType = "COD"
    PaymentTypePrepaid PaymentType = "PREPAID"
)

type NDRAction string
const (
    NDRActionRetry      NDRAction = "RETRY"
    NDRActionReschedule NDRAction = "RESCHEDULE"
    NDRActionRTO        NDRAction = "RTO"
    NDRActionCancel     NDRAction = "CANCEL"
)

// Input Types
type ShippingRateInput struct {
    OriginPincode      int                    `json:"origin_pincode"`
    DestinationPincode int                    `json:"destination_pincode"`
    Weight             float64                `json:"weight"`
    CourierCodes       []string               `json:"courier_codes,omitempty"`
    PaymentMode        PaymentMode            `json:"payment_mode"`
    CollectableAmount  float64                `json:"collectable_amount"`
    Dimensions         []PackageDimensionInput `json:"dimensions,omitempty"`
}

type PackageDimensionInput struct {
    Length float64 `json:"length"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
}

type AvailabilityInput struct {
    OriginPincode      int     `json:"origin_pincode"`
    DestinationPincode int     `json:"destination_pincode"`
    Weight             float64 `json:"weight,omitempty"`
}

type CreateShipmentInput struct {
    CourierCode       string         `json:"courier_code"`
    OrderNumber       string         `json:"order_number"`
    PaymentType       PaymentType    `json:"payment_type"`
    PackageWeight     float64        `json:"package_weight"`
    PackageLength     float64        `json:"package_length"`
    PackageBreadth    float64        `json:"package_breadth"`
    PackageHeight     float64        `json:"package_height"`
    OrderAmount       float64        `json:"order_amount"`
    CollectableAmount float64        `json:"collectable_amount"`
    Consignee        Address        `json:"consignee"`
    Pickup           Address        `json:"pickup"`
    Items            []OrderItem    `json:"items"`
}

type Address struct {
    Name         string `json:"name"`
    CompanyName  string `json:"company_name,omitempty"`
    AddressLine1 string `json:"address_line1"`
    AddressLine2 string `json:"address_line2,omitempty"`
    City         string `json:"city"`
    State        string `json:"state"`
    Pincode      string `json:"pincode"`
    Phone        string `json:"phone"`
    Email        string `json:"email,omitempty"`
    GSTIN        string `json:"gstin,omitempty"`
}

type OrderItem struct {
    Name     string  `json:"name"`
    SKU      string  `json:"sku"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
}

type TrackingInput struct {
    CourierCode string `json:"courier_code"`
    TrackingID  string `json:"tracking_id"`
}

type CancelShipmentInput struct {
    CourierCode string `json:"courier_code"`
    TrackingID  string `json:"tracking_id"`
    Reason      string `json:"reason"`
}

type NDRListInput struct {
    CourierCode string `json:"courier_code"`
    Page        int    `json:"page"`
    Limit       int    `json:"limit"`
}

type NDRActionInput struct {
    CourierCode string    `json:"courier_code"`
    TrackingID  string    `json:"tracking_id"`
    Action      NDRAction `json:"action"`
    Remarks     string    `json:"remarks,omitempty"`
}

// Response Types
type ShippingRateResponse struct {
    Success bool           `json:"success"`
    Rates   []*CourierRate `json:"rates,omitempty"`
    Error   string         `json:"error,omitempty"`
}

type CourierRate struct {
    CourierCode    string       `json:"courier_code"`
    CourierName    string       `json:"courier_name"`
    BaseCharge    float64       `json:"base_charge"`
    FuelSurcharge    float64       `json:"fuel_charge"`
    CodCharge       float64      `json:"cod_charge`
    HandlingCharge float64          `json:"handling_charge"`
    TotalCharge   float64 `json:"total_charge"`
    ExpectedDays    int `json:"expected_days"`
}

type RateDetails struct {
    TotalAmount   float64 `json:"total_amount"`
    BaseAmount    float64 `json:"base_amount"`
    TaxAmount     float64 `json:"tax_amount"`
    CodCharges    float64 `json:"cod_charges"`
    FuelSurcharge float64 `json:"fuel_surcharge"`
}

type CourierAvailabilityResponse struct {
    Success           bool           `json:"success"`
    AvailableCouriers []*CourierInfo `json:"available_couriers,omitempty"`
    Error            string         `json:"error,omitempty"`
}

type CourierInfo struct {
    CourierCode       string   `json:"courier_code"`
    CourierName       string   `json:"courier_name"`
    Description       string   `json:"description,omitempty"`
    SupportedServices []string `json:"supported_services,omitempty"`
}

type ShipmentResponse struct {
    Success    bool   `json:"success"`
    TrackingId string `json:"tracking_id,omitempty"`
    CourierAwb string `json:"courier_awb,omitempty"`
    Label      string `json:"label,omitempty"`
    Error      string `json:"error,omitempty"`
}

type TrackingResponse struct {
    Success bool             `json:"success"`
    Events  []*TrackingEvent `json:"events,omitempty"`
    Error   string           `json:"error,omitempty"`
}

type TrackingEvent struct {
    Status      string `json:"status"`
    Location    string `json:"location"`
    Timestamp   string `json:"timestamp"`
    Description string `json:"description"`
}

type CancelResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error,omitempty"`
}

type NDRListResponse struct {
    Success   bool           `json:"success"`
    Shipments []*NDRShipment `json:"shipments,omitempty"`
    Error     string         `json:"error,omitempty"`
}

type NDRShipment struct {
    TrackingID      string `json:"tracking_id"`
    CourierAWB      string `json:"courier_awb"`
    AttemptCount    int    `json:"attempt_count"`
    LastAttemptDate string `json:"last_attempt_date"`
    Reason          string `json:"reason"`
    Status          string `json:"status"`
}

type NDRActionResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error,omitempty"`
}