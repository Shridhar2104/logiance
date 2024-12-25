package account

import (
    "context"
    "log"
    "time"

    "github.com/Shridhar2104/logilo/account/pb"
    "github.com/google/uuid"
    "google.golang.org/grpc"
)

// Client is a struct that manages the gRPC connection and AccountServiceClient
type Client struct {
    conn    *grpc.ClientConn         // Holds the gRPC client connection
    service pb.AccountServiceClient   // gRPC client for calling the remote AccountService
}

// NewClient establishes a new gRPC connection and returns a Client instance
func NewClient(url string) (*Client, error) {
    // Establish a connection to the gRPC server using the provided URL
    conn, err := grpc.Dial(url, grpc.WithInsecure())
    if err != nil {
        return nil, err // Return error if connection fails
    }

    // Initialize the gRPC AccountService client
    c := pb.NewAccountServiceClient(conn)
    // Return the client with the established connection and gRPC service client
    return &Client{conn: conn, service: c}, nil
}

// Close closes the gRPC connection to release resources
func (c *Client) Close() {
    c.conn.Close()
}

// CreateAccount sends a request to the server to create a new account
func (c *Client) CreateAccount(ctx context.Context, a *Account) (*Account, error) {
    // Send the CreateAccount request to the server
    res, err := c.service.CreateAccount(ctx, &pb.CreateAccountRequest{
        Name:     a.Name,
        Password: a.Password,
        Email:    a.Email,
    })
    
    if err != nil {
        log.Printf("Error creating account: %v", err) // Log the error if the RPC fails
        return nil, err
    }

    // Parse and return the server response as an Account instance
    return &Account{
        ID:        uuid.MustParse(res.Account.Id),    // Parse the account ID from server response
        Name:      res.Account.Name,                  // Map the name field from the server response
        Password:  res.Account.Password,              // Map the password field from server response
        Email:     res.Account.Email,                 // Map the email field from server response
        CreatedAt: time.Now(),                       // Set the current time as creation timestamp
        UpdatedAt: time.Now(),                       // Set the current time as last updated timestamp
    }, nil
}

// GetAccount authenticates and fetches specific account details from the server by email and password
func (c *Client) LoginAndGetAccount(ctx context.Context, email string, password string) (*Account, error) {
    // Send the GetAccountByEmailAndPassword request to the server
    res, err := c.service.GetAccountByEmailAndPassword(ctx, &pb.GetAccountByEmailAndPasswordRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        return nil, err // Return error if RPC fails
    }

    // Parse the server response and map it into an Account instance
    return &Account{
        ID:    uuid.MustParse(res.Account.Id), // Parse the account ID
        Name:  res.Account.Name,               // Map the account name
        Email: email,
    }, nil
}

// ListAccounts fetches a paginated list of accounts from the server
func (c *Client) ListAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
    // Send the ListAccounts request to the server with pagination parameters
    res, err := c.service.ListAccounts(ctx, &pb.ListAccountsRequest{Skip: skip, Take: take})
    if err != nil {
        return nil, err // Handle any RPC failure
    }

    // Map the server response accounts into a slice of Account structs
    accounts := make([]Account, len(res.Accounts)) // Preallocate slice with the correct length
    for i, a := range res.Accounts {
        accounts[i] = Account{
            ID:   uuid.MustParse(a.Id), // Parse and map each account's ID
            Name: a.Name,               // Map the account's name
        }
    }
    return accounts, nil // Return the mapped slice
}

// AddBankAccount sends a request to create a new bank account for a user
func (c *Client) AddBankAccount(ctx context.Context, bankAccount *BankAccount) (*BankAccount, error) {
    // Send the AddBankAccount request to the server
    res, err := c.service.AddBankAccount(ctx, &pb.AddBankAccountRequest{
        UserId:          bankAccount.UserID,
        AccountNumber:   bankAccount.AccountNumber,
        BeneficiaryName: bankAccount.BeneficiaryName,
        IfscCode:        bankAccount.IFSCCode,
        BankName:        bankAccount.BankName,
    })
    
    if err != nil {
        log.Printf("Error adding bank account: %v", err) // Log the error if the RPC fails
        return nil, err
    }

    // Parse and return the server response as a BankAccount instance
    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
        CreatedAt:       time.Now(),                    // Set the current time as creation timestamp
        UpdatedAt:       time.Now(),                    // Set the current time as last updated timestamp
    }, nil
}

// GetBankAccount retrieves bank account details for a specific user
func (c *Client) GetBankAccount(ctx context.Context, userID string) (*BankAccount, error) {
    // Send the GetBankAccount request to the server
    res, err := c.service.GetBankAccount(ctx, &pb.GetBankAccountRequest{
        UserId: userID,
    })
    if err != nil {
        return nil, err // Return error if RPC fails
    }

    // Parse and return the server response as a BankAccount instance
    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
    }, nil
}

// UpdateBankAccount sends a request to update an existing bank account
func (c *Client) UpdateBankAccount(ctx context.Context, bankAccount *BankAccount) (*BankAccount, error) {
    // Send the UpdateBankAccount request to the server
    res, err := c.service.UpdateBankAccount(ctx, &pb.UpdateBankAccountRequest{
        UserId:          bankAccount.UserID,
        AccountNumber:   bankAccount.AccountNumber,
        BeneficiaryName: bankAccount.BeneficiaryName,
        IfscCode:        bankAccount.IFSCCode,
        BankName:        bankAccount.BankName,
    })
    
    if err != nil {
        log.Printf("Error updating bank account: %v", err) // Log the error if the RPC fails
        return nil, err
    }

    // Parse and return the updated bank account details
    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
        UpdatedAt:       time.Now(),                    // Update the last modified timestamp
    }, nil
}

// DeleteBankAccount sends a request to remove a bank account
func (c *Client) DeleteBankAccount(ctx context.Context, userID string) error {
    // Send the DeleteBankAccount request to the server
    _, err := c.service.DeleteBankAccount(ctx, &pb.DeleteBankAccountRequest{
        UserId: userID,
    })
    if err != nil {
        log.Printf("Error deleting bank account: %v", err) // Log the error if the RPC fails
        return err
    }

    return nil
}