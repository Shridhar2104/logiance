package account

import (
    "context"
    "time"

    "github.com/google/uuid"
)

// Extend Service interface with bank account operations
type Service interface {
    CreateAccount(ctx context.Context, name string, password string, email string) (*Account, error)
    LoginAccount(ctx context.Context, email string, password string) (*Account, error)
    ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
    // New bank account methods
    AddBankAccount(ctx context.Context, userID string, accountNumber string, beneficiaryName string, ifscCode string, bankName string) (*BankAccount, error)
    GetBankAccount(ctx context.Context, userID string) (*BankAccount, error)
    UpdateBankAccount(ctx context.Context, userID string, accountNumber string, beneficiaryName string, ifscCode string, bankName string) (*BankAccount, error)
    DeleteBankAccount(ctx context.Context, userID string) error
}

// Keep existing Account struct
type Account struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Password  string    `json:"password"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Add BankAccount struct
type BankAccount struct {
    UserID          string    `json:"user_id"`
    AccountNumber   string    `json:"account_number"`
    BeneficiaryName string    `json:"beneficiary_name"`
    IFSCCode        string    `json:"ifsc_code"`
    BankName        string    `json:"bank_name"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type accountService struct {
    repo Repository
}

func NewAccountService(repo Repository) Service {
    return &accountService{repo}
}

// Keep existing account methods (CreateAccount, LoginAccount, ListAccounts)...
// [Your existing account methods remain unchanged]

// Add new bank account methods
func (s *accountService) AddBankAccount(
    ctx context.Context,
    userID string,
    accountNumber string,
    beneficiaryName string,
    ifscCode string,
    bankName string,
) (*BankAccount, error) {
    bankAccount := &BankAccount{
        UserID:          userID,
        AccountNumber:   accountNumber,
        BeneficiaryName: beneficiaryName,
        IFSCCode:        ifscCode,
        BankName:        bankName,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    err := s.repo.AddBankAccount(ctx, *bankAccount)
    if err != nil {
        return nil, err
    }

    return bankAccount, nil
}

func (s *accountService) GetBankAccount(ctx context.Context, userID string) (*BankAccount, error) {
    return s.repo.GetBankAccount(ctx, userID)
}

func (s *accountService) UpdateBankAccount(
    ctx context.Context,
    userID string,
    accountNumber string,
    beneficiaryName string,
    ifscCode string,
    bankName string,
) (*BankAccount, error) {
    bankAccount := &BankAccount{
        UserID:          userID,
        AccountNumber:   accountNumber,
        BeneficiaryName: beneficiaryName,
        IFSCCode:        ifscCode,
        BankName:        bankName,
        UpdatedAt:       time.Now(),
    }

    err := s.repo.UpdateBankAccount(ctx, *bankAccount)
    if err != nil {
        return nil, err
    }

    return bankAccount, nil
}

func (s *accountService) DeleteBankAccount(ctx context.Context, userID string) error {
    return s.repo.DeleteBankAccount(ctx, userID)
}

// Keep existing CreateAccount method
func (s *accountService) CreateAccount(ctx context.Context, name string, password string, email string) (*Account, error) {
    a := &Account{
        ID:        uuid.New(),
        Name:      name,
        Password:  password,
        Email:     email,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := s.repo.PutAccount(ctx, *a); err != nil {
        return nil, err
    }

    return a, nil
}

// Keep existing LoginAccount method
func (s *accountService) LoginAccount(ctx context.Context, email string, password string) (*Account, error) {
    account, err := s.repo.GetAccountByEmailAndPassword(ctx, email, password)
    if err != nil {
        return nil, err
    }

    return account, nil
}

// Keep existing ListAccounts method
func (s *accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
    if take > 100 || (skip == 0 && take == 0) {
        take = 100
    }

    return s.repo.ListAccounts(ctx, skip, take)
}