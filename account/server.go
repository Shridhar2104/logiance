package account

import (
    "context"
    "fmt"
    "log"
    "net"

    "github.com/Shridhar2104/logilo/account/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

type grpcServer struct {
    pb.UnimplementedAccountServiceServer
    service Service
}

func NewGRPCServer(service Service, port int) error {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return fmt.Errorf("failed to listen on port %d: %w", port, err)
    }

    server := grpc.NewServer()
    pb.RegisterAccountServiceServer(server, &grpcServer{service: service})
    reflection.Register(server)

    log.Printf("gRPC server listening on port %d", port)
    return server.Serve(lis)
}

// Account methods 
func (s *grpcServer) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
    a, err := s.service.CreateAccount(ctx, r.Name, r.Password, r.Email)
    if err != nil {
        log.Printf("Failed to create account: %v", err)
        return nil, fmt.Errorf("failed to create account: %w", err)
    }

    return &pb.CreateAccountResponse{
        Account: &pb.Account{
            Id:   a.ID.String(),
            Name: a.Name,
        },
    }, nil
}

func (s *grpcServer) GetAccountByEmailAndPassword(ctx context.Context, r *pb.GetAccountByEmailAndPasswordRequest) (*pb.GetAccountByEmailAndPasswordResponse, error) {
    a, err := s.service.LoginAccount(ctx, r.Email, r.Password)
    if err != nil {
        log.Printf("Error while authenticating account: %v", err)
        return nil, fmt.Errorf("error while authenticating account: %w", err)
    }

    return &pb.GetAccountByEmailAndPasswordResponse{
        Account: &pb.Account{
            Id:    a.ID.String(),
            Name:  a.Name,
            Email: a.Email,
        },
    }, nil
}

func (s *grpcServer) ListAccounts(ctx context.Context, r *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
    accounts, err := s.service.ListAccounts(ctx, r.Skip, r.Take)
    if err != nil {
        log.Printf("Error while listing accounts: %v", err)
        return nil, fmt.Errorf("error while listing accounts: %w", err)
    }

    grpcAccounts := []*pb.Account{}
    for _, account := range accounts {
        grpcAccounts = append(grpcAccounts, &pb.Account{
            Id:   account.ID.String(),
            Name: account.Name,
        })
    }

    return &pb.ListAccountsResponse{
        Accounts: grpcAccounts,
    }, nil
}

// bank account methods
func (s *grpcServer) AddBankAccount(ctx context.Context, r *pb.AddBankAccountRequest) (*pb.AddBankAccountResponse, error) {
    bankAccount, err := s.service.AddBankAccount(
        ctx,
        r.UserId,
        r.AccountNumber,
        r.AccountType,     // New field
        r.BranchName,      // New field
        r.BeneficiaryName,
        r.IfscCode,
        r.BankName,
    )
    if err != nil {
        log.Printf("Failed to add bank account: %v", err)
        return nil, fmt.Errorf("failed to add bank account: %w", err)
    }

    return &pb.AddBankAccountResponse{
        BankAccount: &pb.BankAccount{
            UserId:          bankAccount.UserID,
            AccountNumber:   bankAccount.AccountNumber,
            AccountType:     bankAccount.AccountType,     // New field
            BranchName:      bankAccount.BranchName,      // New field
            BeneficiaryName: bankAccount.BeneficiaryName,
            IfscCode:        bankAccount.IFSCCode,
            BankName:        bankAccount.BankName,
        },
    }, nil
}

func (s *grpcServer) GetBankAccount(ctx context.Context, r *pb.GetBankAccountRequest) (*pb.GetBankAccountResponse, error) {
    bankAccount, err := s.service.GetBankAccount(ctx, r.UserId)
    if err != nil {
        log.Printf("Error while getting bank account: %v", err)
        return nil, fmt.Errorf("error while getting bank account: %w", err)
    }

    return &pb.GetBankAccountResponse{
        BankAccount: &pb.BankAccount{
            UserId:          bankAccount.UserID,
            AccountNumber:   bankAccount.AccountNumber,
            AccountType:     bankAccount.AccountType,     // New field
            BranchName:      bankAccount.BranchName,      // New field
            BeneficiaryName: bankAccount.BeneficiaryName,
            IfscCode:        bankAccount.IFSCCode,
            BankName:        bankAccount.BankName,
        },
    }, nil
}

func (s *grpcServer) UpdateBankAccount(ctx context.Context, r *pb.UpdateBankAccountRequest) (*pb.UpdateBankAccountResponse, error) {
    bankAccount, err := s.service.UpdateBankAccount(
        ctx,
        r.UserId,
        r.AccountNumber,
        r.AccountType,     // New field
        r.BranchName,      // New field
        r.BeneficiaryName,
        r.IfscCode,
        r.BankName,
    )
    if err != nil {
        log.Printf("Failed to update bank account: %v", err)
        return nil, fmt.Errorf("failed to update bank account: %w", err)
    }

    return &pb.UpdateBankAccountResponse{
        BankAccount: &pb.BankAccount{
            UserId:          bankAccount.UserID,
            AccountNumber:   bankAccount.AccountNumber,
            AccountType:     bankAccount.AccountType,     // New field
            BranchName:      bankAccount.BranchName,      // New field
            BeneficiaryName: bankAccount.BeneficiaryName,
            IfscCode:        bankAccount.IFSCCode,
            BankName:        bankAccount.BankName,
        },
    }, nil
}

func (s *grpcServer) DeleteBankAccount(ctx context.Context, r *pb.DeleteBankAccountRequest) (*pb.DeleteBankAccountResponse, error) {
    err := s.service.DeleteBankAccount(ctx, r.UserId)
    if err != nil {
        log.Printf("Failed to delete bank account: %v", err)
        return nil, fmt.Errorf("failed to delete bank account: %w", err)
    }

    return &pb.DeleteBankAccountResponse{
        Success: true,
    }, nil
}

// Address Method

func (s *grpcServer) AddAddress(ctx context.Context, r *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
    address, err := s.service.AddAddress(
        ctx,
        r.UserId,
        r.ContactPerson,
        r.ContactNumber,
        r.EmailAddress,
        r.CompleteAddress,
        r.Landmark,
        r.Pincode,
        r.City,
        r.State,
        r.Country,
    )
    if err != nil {
        log.Printf("Failed to add address: %v", err)
        return nil, fmt.Errorf("failed to add address: %w", err)
    }

    return &pb.AddAddressResponse{
        Address: &pb.Address{
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
        },
    }, nil
}

func (s *grpcServer) GetAddresses(ctx context.Context, r *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
    addresses, err := s.service.GetAddresses(ctx, r.UserId)
    if err != nil {
        log.Printf("Error while getting addresses: %v", err)
        return nil, fmt.Errorf("error while getting addresses: %w", err)
    }

    var protoAddresses []*pb.Address
    for _, addr := range addresses {
        protoAddresses = append(protoAddresses, &pb.Address{
            Id:              addr.ID.String(),
            UserId:          addr.UserID,
            ContactPerson:   addr.ContactPerson,
            ContactNumber:   addr.ContactNumber,
            EmailAddress:    addr.EmailAddress,
            CompleteAddress: addr.CompleteAddress,
            Landmark:        addr.Landmark,
            Pincode:        addr.Pincode,
            City:           addr.City,
            State:          addr.State,
            Country:        addr.Country,
        })
    }

    return &pb.GetAddressesResponse{
        Addresses: protoAddresses,
    }, nil
}

func (s *grpcServer) UpdateAddress(ctx context.Context, r *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
    address, err := s.service.UpdateAddress(
        ctx,
        r.Id,
        r.UserId,
        r.ContactPerson,
        r.ContactNumber,
        r.EmailAddress,
        r.CompleteAddress,
        r.Landmark,
        r.Pincode,
        r.City,
        r.State,
        r.Country,
    )
    if err != nil {
        log.Printf("Failed to update address: %v", err)
        return nil, fmt.Errorf("failed to update address: %w", err)
    }

    return &pb.UpdateAddressResponse{
        Address: &pb.Address{
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
        },
    }, nil
}

func (s *grpcServer) DeleteAddress(ctx context.Context, r *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
    err := s.service.DeleteAddress(ctx, r.Id)
    if err != nil {
        log.Printf("Failed to delete address: %v", err)
        return nil, fmt.Errorf("failed to delete address: %w", err)
    }

    return &pb.DeleteAddressResponse{
        Success: true,
    }, nil
}

func (s *grpcServer) GetAddressByID(ctx context.Context, r *pb.GetAddressByIDRequest) (*pb.GetAddressByIDResponse, error) {
    address, err := s.service.GetAddressByID(ctx, r.Id)
    if err != nil {
        log.Printf("Error while getting address by ID: %v", err)
        return nil, fmt.Errorf("error while getting address by ID: %w", err)
    }

    return &pb.GetAddressByIDResponse{
        Address: &pb.Address{
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
        },
    }, nil
}