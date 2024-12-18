-- Orders table
ALTER TABLE orders ADD COLUMN delivery_date TIMESTAMP;
ALTER TABLE orders ADD COLUMN remittance_processed BOOLEAN DEFAULT FALSE;

-- Wallets table
CREATE TABLE wallets (
    account_id VARCHAR(255) PRIMARY KEY,
    balance NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table for wallet history
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- e.g., "recharge", "deduction", "remittance"
    amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_id VARCHAR(255) DEFAULT NULL,
    FOREIGN KEY (account_id) REFERENCES wallets(account_id)
);
