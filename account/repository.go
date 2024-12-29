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
    // Address operations
    AddAddress(ctx context.Context, address Address) error
    GetAddresses(ctx context.Context, userID string) ([]Address, error)
    UpdateAddress(ctx context.Context, address Address) error
    DeleteAddress(ctx context.Context, addressID string) error
    GetAddressByID(ctx context.Context, addressID string) (*Address, error)  // Updated return type
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

// AddAddress adds a new address for a user
func (r *postgresRepository) AddAddress(ctx context.Context, address Address) error {
    query := `
        INSERT INTO addresses (id, user_id, contact_person, contact_number, email_address, 
        complete_address, landmark, pincode, city, state, country, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
    
    _, err := r.db.ExecContext(
        ctx,
        query,
        address.ID,
        address.UserID,
        address.ContactPerson,
        address.ContactNumber,
        address.EmailAddress,
        address.CompleteAddress,
        address.Landmark,
        address.Pincode,
        address.City,
        address.State,
        address.Country,
        time.Now(),
        time.Now(),
    )
    
    if err != nil {
        return fmt.Errorf("failed to insert address: %w", err)
    }
    return nil
}

// GetAddresses retrieves all addresses for a user
func (r *postgresRepository) GetAddresses(ctx context.Context, userID string) ([]Address, error) {
    query := `
        SELECT id, user_id, contact_person, contact_number, email_address, 
        complete_address, landmark, pincode, city, state, country, created_at, updated_at
        FROM addresses
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query addresses: %w", err)
    }
    defer rows.Close()

    var addresses []Address
    for rows.Next() {
        var addr Address
        err := rows.Scan(
            &addr.ID,
            &addr.UserID,
            &addr.ContactPerson,
            &addr.ContactNumber,
            &addr.EmailAddress,
            &addr.CompleteAddress,
            &addr.Landmark,
            &addr.Pincode,
            &addr.City,
            &addr.State,
            &addr.Country,
            &addr.CreatedAt,
            &addr.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan address: %w", err)
        }
        addresses = append(addresses, addr)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating address rows: %w", err)
    }
    
    return addresses, nil
}

// UpdateAddress updates an existing address
func (r *postgresRepository) UpdateAddress(ctx context.Context, address Address) error {
    query := `
        UPDATE addresses
        SET contact_person = $3,
            contact_number = $4,
            email_address = $5,
            complete_address = $6,
            landmark = $7,
            pincode = $8,
            city = $9,
            state = $10,
            country = $11,
            updated_at = $12
        WHERE id = $1 AND user_id = $2
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        address.ID,
        address.UserID,
        address.ContactPerson,
        address.ContactNumber,
        address.EmailAddress,
        address.CompleteAddress,
        address.Landmark,
        address.Pincode,
        address.City,
        address.State,
        address.Country,
        time.Now(),
    )
    
    if err != nil {
        return fmt.Errorf("failed to update address: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return errors.New("address not found")
    }
    
    return nil
}

// DeleteAddress removes an address
func (r *postgresRepository) DeleteAddress(ctx context.Context, addressID string) error {
    query := `DELETE FROM addresses WHERE id = $1`
    
    result, err := r.db.ExecContext(ctx, query, addressID)
    if err != nil {
        return fmt.Errorf("failed to delete address: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return errors.New("address not found")
    }
    
    return nil
}

// GetAddressByID retrieves a specific address by its ID
// In repository.go, change the function signature to:
func (r *postgresRepository) GetAddressByID(ctx context.Context, addressID string) (*Address, error) {
    query := `
        SELECT id, user_id, contact_person, contact_number, email_address, 
        complete_address, landmark, pincode, city, state, country, created_at, updated_at
        FROM addresses
        WHERE id = $1
    `
    
    var addr Address
    err := r.db.QueryRowContext(ctx, query, addressID).Scan(
        &addr.ID,
        &addr.UserID,
        &addr.ContactPerson,
        &addr.ContactNumber,
        &addr.EmailAddress,
        &addr.CompleteAddress,
        &addr.Landmark,
        &addr.Pincode,
        &addr.City,
        &addr.State,
        &addr.Country,
        &addr.CreatedAt,
        &addr.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("address not found: %w", err)
        }
        return nil, fmt.Errorf("failed to query address: %w", err)
    }
    
    return &addr, nil
}