-- accounts/up.sql

-- Create the accounts table (your existing table)
DROP TABLE IF EXISTS accounts;
CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create the bank_accounts table
DROP TABLE IF EXISTS bank_accounts;
CREATE TABLE IF NOT EXISTS bank_accounts (
    user_id VARCHAR(36) PRIMARY KEY,
    account_number VARCHAR(50) NOT NULL,
    beneficiary_name VARCHAR(255) NOT NULL,
    ifsc_code VARCHAR(11) NOT NULL,  -- IFSC codes are typically 11 characters
    bank_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES accounts(id) ON DELETE CASCADE
);