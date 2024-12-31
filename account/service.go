package account

import (
    "context"
    "time"
    "fmt"
    "github.com/google/uuid"
)

// Extend Service interface with bank account operations
type Service interface {
    CreateAccount(ctx context.Context, name string, password string, email string) (*Account, error)
    LoginAccount(ctx context.Context, email string, password string) (*Account, error)
    ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
    // New bank account methods
    AddBankAccount(ctx context.Context, userID string, accountNumber string, 
        accountType string, branchName string, beneficiaryName string, 
        ifscCode string, bankName string) (*BankAccount, error)
    GetBankAccount(ctx context.Context, userID string) (*BankAccount, error)
    UpdateBankAccount(ctx context.Context, userID string, accountNumber string, 
        accountType string, branchName string, beneficiaryName string, 
        ifscCode string, bankName string) (*BankAccount, error)
    DeleteBankAccount(ctx context.Context, userID string) error
    // Address operations
    AddAddress(ctx context.Context, userID string, contactPerson string, contactNumber string, 
        emailAddress string, completeAddress string, landmark string, pincode string, 
        city string, state string, country string) (*Address, error)
    GetAddresses(ctx context.Context, userID string) ([]Address, error)
    UpdateAddress(ctx context.Context, addressID string, userID string, contactPerson string, 
        contactNumber string, emailAddress string, completeAddress string, landmark string, 
        pincode string, city string, state string, country string) (*Address, error)
    DeleteAddress(ctx context.Context, addressID string) error
    GetAddressByID(ctx context.Context, addressID string) (*Address, error)
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
    AccountType     string    `json:"account_type"`
    BranchName      string    `json:"branch_name"`
    BeneficiaryName string    `json:"beneficiary_name"`
    IFSCCode        string    `json:"ifsc_code"`
    BankName        string    `json:"bank_name"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
// Address represents the user address data structure
type Address struct {
    ID             uuid.UUID     `json:"id"`
    UserID         string    `json:"user_id"`
    ContactPerson  string    `json:"contact_person"`
    ContactNumber  string    `json:"contact_number"`
    EmailAddress   string    `json:"email_address"`
    CompleteAddress string   `json:"complete_address"`
    Landmark       string    `json:"landmark"`
    Pincode        string    `json:"pincode"`
    City           string    `json:"city"`
    State          string    `json:"state"`
    Country        string    `json:"country"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type accountService struct {
    repo Repository
}

func NewAccountService(repo Repository) Service {
    return &accountService{repo}
}


// Bank account methods
func (s *accountService) AddBankAccount(
    ctx context.Context,
    userID string,
    accountNumber string,
    accountType string,
    branchName string,
    beneficiaryName string,
    ifscCode string,
    bankName string,
) (*BankAccount, error) {
    bankAccount := &BankAccount{
        UserID:          userID,
        AccountNumber:   accountNumber,
        AccountType:     accountType,
        BranchName:      branchName,
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
    accountType string,
    branchName string,
    beneficiaryName string,
    ifscCode string,
    bankName string,
) (*BankAccount, error) {
    bankAccount := &BankAccount{
        UserID:          userID,
        AccountNumber:   accountNumber,
        AccountType:     accountType,
        BranchName:      branchName,
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

// CreateAccount method
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

//  LoginAccount method
func (s *accountService) LoginAccount(ctx context.Context, email string, password string) (*Account, error) {
    account, err := s.repo.GetAccountByEmailAndPassword(ctx, email, password)
    if err != nil {
        return nil, err
    }

    return account, nil
}

//  ListAccounts method
func (s *accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
    if take > 100 || (skip == 0 && take == 0) {
        take = 100
    }

    return s.repo.ListAccounts(ctx, skip, take)
}
//method implementations to  accountService struct

func (s *accountService) AddAddress(
    ctx context.Context,
    userID string,
    contactPerson string,
    contactNumber string,
    emailAddress string,
    completeAddress string,
    landmark string,
    pincode string,
    city string,
    state string,
    country string,
) (*Address, error) {
    address := &Address{
        ID:             uuid.New(),
        UserID:         userID,
        ContactPerson:  contactPerson,
        ContactNumber:  contactNumber,
        EmailAddress:   emailAddress,
        CompleteAddress: completeAddress,
        Landmark:       landmark,
        Pincode:        pincode,
        City:           city,
        State:          state,
        Country:        country,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    err := s.repo.AddAddress(ctx, *address)
    if err != nil {
        return nil, err
    }

    return address, nil
}

func (s *accountService) GetAddresses(ctx context.Context, userID string) ([]Address, error) {
    return s.repo.GetAddresses(ctx, userID)
}

func (s *accountService) UpdateAddress(
    ctx context.Context,
    addressID string,  // We'll keep this as string in the parameter
    userID string,
    contactPerson string,
    contactNumber string,
    emailAddress string,
    completeAddress string,
    landmark string,
    pincode string,
    city string,
    state string,
    country string,
) (*Address, error) {
    // Parse the string ID to UUID
    addressUUID, err := uuid.Parse(addressID)
    if err != nil {
        return nil, fmt.Errorf("invalid address ID: %w", err)
    }

    address := &Address{
        ID:             addressUUID,  // Now using UUID type
        UserID:         userID,
        ContactPerson:  contactPerson,
        ContactNumber:  contactNumber,
        EmailAddress:   emailAddress,
        CompleteAddress: completeAddress,
        Landmark:       landmark,
        Pincode:        pincode,
        City:           city,
        State:          state,
        Country:        country,
        UpdatedAt:      time.Now(),
    }

    err = s.repo.UpdateAddress(ctx, *address)
    if err != nil {
        return nil, err
    }

    return address, nil
}
func (s *accountService) DeleteAddress(ctx context.Context, addressID string) error {
    return s.repo.DeleteAddress(ctx, addressID)
}

func (s *accountService) GetAddressByID(ctx context.Context, addressID string) (*Address, error) {
    return s.repo.GetAddressByID(ctx, addressID)
}