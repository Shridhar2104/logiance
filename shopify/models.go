package shopify

import (
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)


// Order represents a Shopify order and is structured for storage and processing.
type Order struct {
	ID                int64            `json:"id" db:"id"` // Primary key
	Name              string           `json:"name,omitempty" db:"name"`
	Email             string           `json:"email,omitempty" db:"email"`
	CreatedAt         time.Time        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at,omitempty" db:"updated_at"`
	CancelledAt       *time.Time       `json:"cancelled_at,omitempty" db:"cancelled_at"`
	ClosedAt          *time.Time       `json:"closed_at,omitempty" db:"closed_at"`
	ProcessedAt       *time.Time       `json:"processed_at,omitempty" db:"processed_at"`
	Currency          string           `json:"currency,omitempty" db:"currency"`
	TotalPrice        float64          `json:"total_price,omitempty" db:"total_price"`
	SubtotalPrice     float64          `json:"subtotal_price,omitempty" db:"subtotal_price"`
	TotalDiscounts    float64          `json:"total_discounts,omitempty" db:"total_discounts"`
	TotalTax          float64          `json:"total_tax,omitempty" db:"total_tax"`
	TaxesIncluded     bool             `json:"taxes_included,omitempty" db:"taxes_included"`
	TotalWeight       int              `json:"total_weight,omitempty" db:"total_weight"`
	FinancialStatus   string           `json:"financial_status,omitempty" db:"financial_status"`
	FulfillmentStatus string           `json:"fulfillment_status,omitempty" db:"fulfillment_status"`
	Number            int              `json:"number,omitempty" db:"number"`
	OrderNumber       int              `json:"order_number,omitempty" db:"order_number"`
	Test              bool             `json:"test,omitempty" db:"test"`
	BrowserIp         string           `json:"browser_ip,omitempty" db:"browser_ip"`
	CancelReason      string           `json:"cancel_reason,omitempty" db:"cancel_reason"`
	Tags              string           `json:"tags,omitempty" db:"tags"`
	LocationId        int64            `json:"location_id,omitempty" db:"location_id"`
	Customer          goshopify.Customer `json:"customer,omitempty" db:"customer"`
	BillingAddress    goshopify.Address  `json:"billing_address,omitempty" db:"billing_address"`
	ShippingAddress   goshopify.Address  `json:"shipping_address,omitempty" db:"shipping_address"`
	Gateway           string           `json:"gateway,omitempty" db:"gateway"`
	Confirmed         bool             `json:"confirmed,omitempty" db:"confirmed"`
	Phone             string           `json:"phone,omitempty" db:"phone"`
	ContactEmail      string           `json:"contact_email,omitempty" db:"contact_email"`
}




type ShopifyOrder struct{
	Order *goshopify.Order
}
