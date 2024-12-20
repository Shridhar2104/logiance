package shopify

import (
	"context"
	"fmt"

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