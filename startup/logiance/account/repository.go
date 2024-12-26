package account

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"
    "golang.org/x/crypto/bcrypt"

    _ "github.com/lib/pq"
)

// // Account represents the account data structure
// type Account struct {
//     ID        string
//     Name      string
//     Email     string
//     Password  string
//     CreatedAt time.Time
//     UpdatedAt time.Time
// }

// // BankAccount represents the bank account data structure
// type BankAccount struct {
//     UserID          string
//     AccountNumber   string
//     BeneficiaryName string
//     IFSCCode        string
//     BankName        string
//     CreatedAt       time.Time
//     UpdatedAt       time.Time
// }

// Repository defines the interface for interacting with the accounts database
type Repository interface {
    Close()
    PutAccount(ctx context.Context, account Account) error
    GetAccountByEmailAndPassword(ctx context.Context, email, password string) (*Account, error)
    ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
    Ping() error
    // Bank account operations
    AddBankAccount(ctx context.Context, bankAccount BankAccount) error
    GetBankAccount(ctx context.Context, userID string) (*BankAccount, error)
    UpdateBankAccount(ctx context.Context, bankAccount BankAccount) error
    DeleteBankAccount(ctx context.Context, userID string) error
}

// postgresRepository is the PostgreSQL implementation of the Repository interface
type postgresRepository struct {
    db *sql.DB
}

// NewPostgresRepository creates and initializes a new postgresRepository instance
func NewPostgresRepository(url string) (Repository, error) {
    db, err := sql.Open("postgres", url)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return &postgresRepository{db}, nil
}

// Close releases the database connection resources
func (r *postgresRepository) Close() {
    r.db.Close()
}

// Ping checks the health of the database connection
func (r *postgresRepository) Ping() error {
    return r.db.Ping()
}

// PutAccount inserts a new account into the accounts table
func (r *postgresRepository) PutAccount(ctx context.Context, account Account) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    query := `
        INSERT INTO accounts (id, name, email, password, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err = r.db.ExecContext(ctx, query, account.ID, account.Name, account.Email, string(hashedPassword), account.CreatedAt, account.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }
    return nil
}

// GetAccountByEmailAndPassword retrieves an account by email and validates password
func (r *postgresRepository) GetAccountByEmailAndPassword(ctx context.Context, email, password string) (*Account, error) {
    query := `
        SELECT id, name, email, password, created_at, updated_at 
        FROM accounts 
        WHERE email = $1
    `
    row := r.db.QueryRowContext(ctx, query, email)

    var account Account
    if err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Password, &account.CreatedAt, &account.UpdatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("email not found: %w", err)
        }
        return nil, fmt.Errorf("failed to query account: %w", err)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid password")
    }

    return &account, nil
}

// ListAccounts retrieves a paginated list of accounts
func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
    query := `
        SELECT id, name, email 
        FROM accounts 
        ORDER BY id DESC 
        LIMIT $1 OFFSET $2
    `
    rows, err := r.db.QueryContext(ctx, query, take, skip)
    if err != nil {
        return nil, fmt.Errorf("failed to query accounts: %w", err)
    }
    defer rows.Close()

    accounts := []Account{}
    for rows.Next() {
        var a Account
        if err := rows.Scan(&a.ID, &a.Name, &a.Email); err != nil {
            return nil, fmt.Errorf("failed to scan account: %w", err)
        }
        accounts = append(accounts, a)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }
    return accounts, nil
}

// AddBankAccount adds a new bank account for a user
func (r *postgresRepository) AddBankAccount(ctx context.Context, bankAccount BankAccount) error {
    query := `
        INSERT INTO bank_accounts (user_id, account_number, beneficiary_name, ifsc_code, bank_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.ExecContext(
        ctx,
        query,
        bankAccount.UserID,
        bankAccount.AccountNumber,
        bankAccount.BeneficiaryName,
        bankAccount.IFSCCode,
        bankAccount.BankName,
        time.Now(),
        time.Now(),
    )
    
    if err != nil {
        return fmt.Errorf("failed to insert bank account: %w", err)
    }
    return nil
}

// GetBankAccount retrieves bank account details for a user
func (r *postgresRepository) GetBankAccount(ctx context.Context, userID string) (*BankAccount, error) {
    query := `
        SELECT user_id, account_number, beneficiary_name, ifsc_code, bank_name, created_at, updated_at
        FROM bank_accounts
        WHERE user_id = $1
    `
    
    var bankAccount BankAccount
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &bankAccount.UserID,
        &bankAccount.AccountNumber,
        &bankAccount.BeneficiaryName,
        &bankAccount.IFSCCode,
        &bankAccount.BankName,
        &bankAccount.CreatedAt,
        &bankAccount.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("bank account not found for user: %w", err)
        }
        return nil, fmt.Errorf("failed to query bank account: %w", err)
    }
    
    return &bankAccount, nil
}

// UpdateBankAccount updates an existing bank account
func (r *postgresRepository) UpdateBankAccount(ctx context.Context, bankAccount BankAccount) error {
    query := `
        UPDATE bank_accounts
        SET account_number = $2,
            beneficiary_name = $3,
            ifsc_code = $4,
            bank_name = $5,
            updated_at = $6
        WHERE user_id = $1
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        bankAccount.UserID,
        bankAccount.AccountNumber,
        bankAccount.BeneficiaryName,
        bankAccount.IFSCCode,
        bankAccount.BankName,
        time.Now(),
    )
    
    if err != nil {
        return fmt.Errorf("failed to update bank account: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return errors.New("bank account not found")
    }
    
    return nil
}

// DeleteBankAccount removes a bank account
func (r *postgresRepository) DeleteBankAccount(ctx context.Context, userID string) error {
    query := `DELETE FROM bank_accounts WHERE user_id = $1`
    
    result, err := r.db.ExecContext(ctx, query, userID)
    if err != nil {
        return fmt.Errorf("failed to delete bank account: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return errors.New("bank account not found")
    }
    
    return nil
}