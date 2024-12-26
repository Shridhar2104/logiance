package shopify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"net"

	"github.com/Shridhar2104/logilo/shopify/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
type grpcServer struct {
	pb.UnimplementedShopifyServiceServer
	service Service
}

func NewGRPCServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterShopifyServiceServer(server, &grpcServer{
		UnimplementedShopifyServiceServer: pb.UnimplementedShopifyServiceServer{}, // Add this
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func(s * grpcServer) GetAuthorizationURL(ctx context.Context, r *pb.GetAuthorizationURLRequest) (*pb.GetAuthorizationURLResponse, error){
	authUrl , err:= s.service.GenerateAuthURL(ctx,r.ShopName, r.State)

	if err!= nil{
		return nil, err
	}

	return &pb.GetAuthorizationURLResponse{
		AuthUrl: authUrl,
		}, nil
}

func(s *grpcServer) ExchangeAccessToken(ctx context.Context, r *pb.ExchangeAccessTokenRequest) (*pb.ExchangeAccessTokenResponse, error){
	err := s.service.ExchangeAccessToken(ctx,r.ShopName,r.Code, r.AccountId)

	if err!= nil{
		return &pb.ExchangeAccessTokenResponse{
			Success: false,
		}, err
	
	}
	return &pb.ExchangeAccessTokenResponse{
		Success: true,
	}, nil

}

// server.go
func (s *grpcServer) SyncOrders(ctx context.Context, r *pb.SyncOrderRequest) (*pb.SyncOrderResponse, error) {
    results, err := s.service.SyncOrdersForAccount(ctx, r.AccountId)
    if err != nil {
        return &pb.SyncOrderResponse{
            OverallSuccess: false,
        }, fmt.Errorf("failed to sync orders: %w", err)
    }

    return &pb.SyncOrderResponse{
        OverallSuccess: true,
        ShopResults: results,
    }, nil
}

// Update the GetOrdersForShopAndAccount function in server.go
func (s *grpcServer) GetOrdersForShopAndAccount(ctx context.Context, r *pb.GetOrdersForShopAndAccountRequest) (*pb.GetOrdersForShopAndAccountResponse, error) {
    // Now we only need accountId, as we'll fetch orders from all shops
    orders, _, err := s.service.GetAccountOrders(ctx, r.AccountId, 1, 100) // You can make page and pageSize configurable
    if err != nil {
        return nil, fmt.Errorf("failed to get orders: %w", err)
    }

    ordersPb := make([]*pb.Order, len(orders))
    for i, order := range orders {
        ordersPb[i] = &pb.Order{
            Id:         fmt.Sprintf("%d", order.ID),
            AccountId:  r.AccountId,
            TotalPrice: float32(order.TotalPrice),
            OrderName:    order.Name,
        }
    }

    return &pb.GetOrdersForShopAndAccountResponse{
        Orders: ordersPb,
    }, nil
}
// func (s *grpcServer) GetOrdersForShopAndAccount(ctx context.Context, r *pb.GetOrdersForShopAndAccountRequest) (*pb.GetOrdersForShopAndAccountResponse, error){
// 	orders, err := s.service.GetOrdersForShopAndAccount(ctx, r.ShopName, r.AccountId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ordersPb := make([]*pb.Order, len(orders))
// 	for i, order := range orders {
// 		ordersPb[i] = &pb.Order{
// 			Id: order.ID,
// 			AccountId: order.AccountId,
// 			ShopId: order.ShopName,
// 			TotalPrice: float32(order.TotalPrice),
// 			OrderId: order.OrderId,
// 		}
// 	}
// 	return &pb.GetOrdersForShopAndAccountResponse{
// 		Orders: ordersPb,
// 	}, nil

// }

func (s *grpcServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
    order, err := s.service.GetOrder(ctx, req.OrderId)
    if err != nil {
        return nil, err
    }

    return &pb.Order{
        Id:               string(order.ID),
        OrderName:         order.Name,
        TotalPrice:        float32(order.TotalPrice),
        SubtotalPrice:     float32(order.SubtotalPrice),
        TotalTax:          float32(order.TotalTax),
        Currency:          order.Currency,
        FinancialStatus:   order.FinancialStatus,
        FulfillmentStatus: order.FulfillmentStatus,
        CreatedAt:         order.CreatedAt.Format(time.RFC3339),
        CustomerEmail:     order.Customer.Email,
        CustomerName:      fmt.Sprintf("%s %s", order.Customer.FirstName, order.Customer.LastName),
    }, nil
}
// Add this function to server.go
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
    if req.PageSize == 0 {
        req.PageSize = 20 // default page size
    }
    if req.Page == 0 {
        req.Page = 1 // default page
    }

    orders, totalCount, err := s.service.GetAccountOrders(
        ctx,
        req.AccountId,
        int(req.Page),
        int(req.PageSize),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get orders: %w", err)
    }

    // Convert orders to proto format
    pbOrders := make([]*pb.Order, len(orders))
    for i, order := range orders {
        customerName := ""
        if order.Customer.FirstName != "" || order.Customer.LastName != "" {
            customerName = strings.TrimSpace(order.Customer.FirstName + " " + order.Customer.LastName)
        }

        pbOrders[i] = &pb.Order{
            Id:                fmt.Sprintf("%d", order.ID),
            OrderName:         order.Name,
            AccountId:         req.AccountId,
            TotalPrice:        float32(order.TotalPrice),
            SubtotalPrice:     float32(order.SubtotalPrice),
            TotalTax:         float32(order.TotalTax),
            Currency:         order.Currency,
            FinancialStatus:  order.FinancialStatus,
            FulfillmentStatus: order.FulfillmentStatus,
            CreatedAt:        order.CreatedAt.Format(time.RFC3339),
            CustomerEmail:    order.Customer.Email,
            CustomerName:     customerName,
        }
    }

    // Calculate total pages
    pageSize := int(req.PageSize)
    totalPages := (totalCount + pageSize - 1) / pageSize

    return &pb.GetOrdersForAccountResponse{
        Orders:      pbOrders,
        TotalCount:  int32(totalCount),
        TotalPages:  int32(totalPages),
        CurrentPage: req.Page,
    }, nil
}