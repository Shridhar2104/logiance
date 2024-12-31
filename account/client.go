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
    res, err := c.service.AddBankAccount(ctx, &pb.AddBankAccountRequest{
        UserId:          bankAccount.UserID,
        AccountNumber:   bankAccount.AccountNumber,
        AccountType:     bankAccount.AccountType,     // New field
        BranchName:      bankAccount.BranchName,      // New field
        BeneficiaryName: bankAccount.BeneficiaryName,
        IfscCode:        bankAccount.IFSCCode,
        BankName:        bankAccount.BankName,
    })
    
    if err != nil {
        log.Printf("Error adding bank account: %v", err)
        return nil, err
    }

    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        AccountType:     res.BankAccount.AccountType,     // New field
        BranchName:      res.BankAccount.BranchName,      // New field
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }, nil
}

// GetBankAccount retrieves bank account details for a specific user
func (c *Client) GetBankAccount(ctx context.Context, userID string) (*BankAccount, error) {
    res, err := c.service.GetBankAccount(ctx, &pb.GetBankAccountRequest{
        UserId: userID,
    })
    if err != nil {
        return nil, err
    }

    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        AccountType:     res.BankAccount.AccountType,     // New field
        BranchName:      res.BankAccount.BranchName,      // New field
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
    }, nil
}

// UpdateBankAccount sends a request to update an existing bank account
func (c *Client) UpdateBankAccount(ctx context.Context, bankAccount *BankAccount) (*BankAccount, error) {
    res, err := c.service.UpdateBankAccount(ctx, &pb.UpdateBankAccountRequest{
        UserId:          bankAccount.UserID,
        AccountNumber:   bankAccount.AccountNumber,
        AccountType:     bankAccount.AccountType,     // New field
        BranchName:      bankAccount.BranchName,      // New field
        BeneficiaryName: bankAccount.BeneficiaryName,
        IfscCode:        bankAccount.IFSCCode,
        BankName:        bankAccount.BankName,
    })
    
    if err != nil {
        log.Printf("Error updating bank account: %v", err)
        return nil, err
    }

    return &BankAccount{
        UserID:          res.BankAccount.UserId,
        AccountNumber:   res.BankAccount.AccountNumber,
        AccountType:     res.BankAccount.AccountType,     // New field
        BranchName:      res.BankAccount.BranchName,      // New field
        BeneficiaryName: res.BankAccount.BeneficiaryName,
        IFSCCode:        res.BankAccount.IfscCode,
        BankName:        res.BankAccount.BankName,
        UpdatedAt:       time.Now(),
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

// AddAddress sends a request to create a new address for a user
func (c *Client) AddAddress(ctx context.Context, address *Address) (*Address, error) {
    res, err := c.service.AddAddress(ctx, &pb.AddAddressRequest{
        UserId:          address.UserID,
        ContactPerson:   address.ContactPerson,
        ContactNumber:   address.ContactNumber,
        EmailAddress:    address.EmailAddress,
        CompleteAddress: address.CompleteAddress,
        Landmark:        address.Landmark,
        Pincode:        address.Pincode,
        City:           address.City,
        State:          address.State,
        Country:        address.Country,
    })
    
    if err != nil {
        log.Printf("Error adding address: %v", err)
        return nil, err
    }

    return &Address{
        ID:              uuid.MustParse(res.Address.Id),
        UserID:          res.Address.UserId,
        ContactPerson:   res.Address.ContactPerson,
        ContactNumber:   res.Address.ContactNumber,
        EmailAddress:    res.Address.EmailAddress,
        CompleteAddress: res.Address.CompleteAddress,
        Landmark:        res.Address.Landmark,
        Pincode:        res.Address.Pincode,
        City:           res.Address.City,
        State:          res.Address.State,
        Country:        res.Address.Country,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }, nil
}

// GetAddresses retrieves all addresses for a specific user
func (c *Client) GetAddresses(ctx context.Context, userID string) ([]Address, error) {
    res, err := c.service.GetAddresses(ctx, &pb.GetAddressesRequest{
        UserId: userID,
    })
    if err != nil {
        return nil, err
    }

    addresses := make([]Address, len(res.Addresses))
    for i, addr := range res.Addresses {
        addresses[i] = Address{
            ID:              uuid.MustParse(addr.Id),
            UserID:          addr.UserId,
            ContactPerson:   addr.ContactPerson,
            ContactNumber:   addr.ContactNumber,
            EmailAddress:    addr.EmailAddress,
            CompleteAddress: addr.CompleteAddress,
            Landmark:        addr.Landmark,
            Pincode:        addr.Pincode,
            City:           addr.City,
            State:          addr.State,
            Country:        addr.Country,
        }
    }
    return addresses, nil
}

// UpdateAddress sends a request to update an existing address
func (c *Client) UpdateAddress(ctx context.Context, address *Address) (*Address, error) {
    res, err := c.service.UpdateAddress(ctx, &pb.UpdateAddressRequest{
        Id:              address.ID.String(),
        UserId:          address.UserID,
        ContactPerson:   address.ContactPerson,
        ContactNumber:   address.ContactNumber,
        EmailAddress:    address.EmailAddress,
        CompleteAddress: address.CompleteAddress,
        Landmark:        address.Landmark,
        Pincode:        address.Pincode,
        City:           address.City,
        State:          address.State,
        Country:        address.Country,
    })
    
    if err != nil {
        log.Printf("Error updating address: %v", err)
        return nil, err
    }

    return &Address{
        ID:              uuid.MustParse(res.Address.Id),
        UserID:          res.Address.UserId,
        ContactPerson:   res.Address.ContactPerson,
        ContactNumber:   res.Address.ContactNumber,
        EmailAddress:    res.Address.EmailAddress,
        CompleteAddress: res.Address.CompleteAddress,
        Landmark:        res.Address.Landmark,
        Pincode:        res.Address.Pincode,
        City:           res.Address.City,
        State:          res.Address.State,
        Country:        res.Address.Country,
        UpdatedAt:       time.Now(),
    }, nil
}

// DeleteAddress sends a request to remove an address
func (c *Client) DeleteAddress(ctx context.Context, addressID string) error {
    _, err := c.service.DeleteAddress(ctx, &pb.DeleteAddressRequest{
        Id: addressID,
    })
    if err != nil {
        log.Printf("Error deleting address: %v", err)
        return err
    }

    return nil
}

// GetAddressByID retrieves a specific address by its ID
func (c *Client) GetAddressByID(ctx context.Context, addressID string) (*Address, error) {
    res, err := c.service.GetAddressByID(ctx, &pb.GetAddressByIDRequest{
        Id: addressID,
    })
    if err != nil {
        log.Printf("Error getting address by ID: %v", err)
        return nil, err
    }

    return &Address{
        ID:              uuid.MustParse(res.Address.Id),
        UserID:          res.Address.UserId,
        ContactPerson:   res.Address.ContactPerson,
        ContactNumber:   res.Address.ContactNumber,
        EmailAddress:    res.Address.EmailAddress,
        CompleteAddress: res.Address.CompleteAddress,
        Landmark:        res.Address.Landmark,
        Pincode:        res.Address.Pincode,
        City:           res.Address.City,
        State:          res.Address.State,
        Country:        res.Address.Country,
    }, nil
}