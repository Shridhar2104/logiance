package models

import "time"

type Account struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Email string `json:"email"`
	Orders []Order `json:"orders"`
	ShopNames []ShopName `json:"shopnames"`
}

type ShopName struct {
	Shopname string `json:"shopname"`
}

type Order struct {
	ID string `json:"id"`
	Amount float64 `json:"amount"`
	AccountID string `json:"accountId"`
	CreatedAt time.Time `json:"createdAt"`
	Description string `json:"description"`
	LineItems []OrderLineItem `json:"lineItems"`
}

type OrderLineItem struct {
	ID string `json:"id"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
}

type AccountInput struct {
	Name string `json:"name"`
}
type OrderInput struct {
	AccountID string `json:"accountId"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
}

type ShopSyncStatus struct {
    Success      bool    `json:"success"`
    ErrorMessage string  `json:"errorMessage,omitempty"`
    OrdersCount  int     `json:"ordersCount"`
}

type ShopSyncDetails struct {
    ShopName string         `json:"shopName"`
    Status   *ShopSyncStatus `json:"status"`
}

type SyncOrdersResult struct {
    OverallSuccess bool               `json:"overallSuccess"`
    Message        string            `json:"message,omitempty"`
    ShopResults    []*ShopSyncDetails `json:"shopResults"`
}



type ShippingRateInput struct {
    OriginPincode      int                    `json:"originPincode"`
    DestinationPincode int                    `json:"destinationPincode"`
    Weight             int                    `json:"weight"`
    CourierCodes       []string               `json:"courierCodes,omitempty"`
    PaymentMode        PaymentMode            `json:"paymentMode"`
    Dimensions         []PackageDimensionInput `json:"dimensions,omitempty"`
}

type PackageDimensionInput struct {
    Length float64 `json:"length"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
}

type AvailabilityInput struct {
    OriginPincode      int `json:"originPincode"`
    DestinationPincode int `json:"destinationPincode"`
}

type PaymentMode string

const (
    PaymentModeCOD     PaymentMode = "COD"
    PaymentModePREPAID PaymentMode = "PREPAID"
)

type ShippingRateResponse struct {
    Success bool           `json:"success"`
    Rates   []*CourierRate `json:"rates,omitempty"`
    Error   string        `json:"error,omitempty"`
}

type CourierRate struct {
    CourierCode    string       `json:"courierCode"`
    CourierName    string       `json:"courierName"`
    ServiceType    string       `json:"serviceType"`
    RateDetails    *RateDetails `json:"rateDetails"`
    EstimatedDays  int         `json:"estimatedDays"`
}

type RateDetails struct {
    TotalAmount    float64 `json:"totalAmount"`
    GrossAmount    float64 `json:"grossAmount"`
    TaxAmount      float64 `json:"taxAmount"`
    CodCharges     float64 `json:"codCharges"`
    FuelSurcharge  float64 `json:"fuelSurcharge"`
}

type CourierAvailabilityResponse struct {
    Success           bool           `json:"success"`
    AvailableCouriers []*CourierInfo `json:"availableCouriers"`
    Error            string        `json:"error,omitempty"`
}

type CourierInfo struct {
    CourierCode       string   `json:"courierCode"`
    CourierName       string   `json:"courierName"`
    SupportedServices []string `json:"supportedServices"`
}