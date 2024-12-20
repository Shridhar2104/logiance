CREATE TABLE tokens (
    shop_name TEXT NOT NULL,
    account_id TEXT NOT NULL,
    token TEXT NOT NULL,
    PRIMARY KEY (shop_name, account_id)
);

DROP TABLE IF EXISTS orders;
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    name TEXT,
    email TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    cancelled_at TIMESTAMP,
    closed_at TIMESTAMP,
    processed_at TIMESTAMP,
    currency TEXT,
    total_price NUMERIC(10,2),
    subtotal_price NUMERIC(10,2),
    total_discounts NUMERIC(10,2),
    total_tax NUMERIC(10,2),
    taxes_included BOOLEAN,
    financial_status TEXT,
    fulfillment_status TEXT,
    order_number INTEGER,
    test BOOLEAN,
    browser_ip TEXT,
    cancel_reason TEXT,
    tags TEXT,
    gateway TEXT,
    confirmed BOOLEAN,
    phone TEXT,
    contact_email TEXT,
    
    -- Additional fields to track ownership
    shop_name TEXT NOT NULL,
    account_id TEXT NOT NULL,
    
    -- Customer fields
    customer_id BIGINT,
    customer_email TEXT,
    customer_first_name TEXT,
    customer_last_name TEXT,
    customer_phone TEXT
);

-- Add indexes for common queries
CREATE INDEX idx_orders_shop_account ON orders(shop_name, account_id);
CREATE INDEX idx_orders_crea