CREATE TABLE tokens (
    shop_name TEXT NOT NULL,
    account_id TEXT NOT NULL,
    token TEXT NOT NULL,
    PRIMARY KEY (shop_name, account_id)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    shop_name TEXT NOT NULL,
    account_id TEXT NOT NULL,
    order_id TEXT NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL
);
