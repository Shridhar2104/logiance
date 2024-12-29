// graphql/graph.go
package main

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Shridhar2104/logilo/account"
	"github.com/Shridhar2104/logilo/graphql/models"
	"github.com/Shridhar2104/logilo/payment"
	"github.com/Shridhar2104/logilo/shopify"

	pb "github.com/Shridhar2104/logilo/shipment/proto/proto"

	"google.golang.org/grpc"
)
type Server struct {
    accountClient *account.Client
    shopifyClient *shopify.Client
    paymentClient *payment.Client
    shipmentClient pb.ShipmentServiceClient
    shippingResolver *ShippingResolver

}



func NewGraphQLServer(accountUrl, shopifyUrl, shipmentUrl, paymentUrl string) (*Server, error) {
    accountClient, err := account.NewClient(accountUrl)
    if err != nil {
        return nil, err
    }

    shopifyClient, err := shopify.NewClient(shopifyUrl)
    if err != nil {
        accountClient.Close()
        return nil, err
    }

    paymentClient, err := payment.NewClient(paymentUrl)
    if err != nil {
        accountClient.Close()
        shopifyClient.Close()
        return nil, err
    }

    // Connect to shipment service
    shipmentConn, err := grpc.Dial(shipmentUrl, grpc.WithInsecure())
    if err != nil {
        log.Printf("Failed to connect to shipment service: %v", err)
        accountClient.Close()
        return nil, err
    }

    shipmentClient := pb.NewShipmentServiceClient(shipmentConn)

        
    server := &Server{
        accountClient:  accountClient,
        shopifyClient:  shopifyClient,
        shipmentClient: shipmentClient,
        paymentClient: paymentClient,
    

    }
    server.shippingResolver = NewShippingResolver(shipmentClient)
    
    // Verify client initialization
    log.Printf("Server initialized with clients - Account: %v, Shopify: %v, Shipment: %v", 
        accountClient != nil, 
        shopifyClient != nil, 
        shipmentClient != nil,
        paymentClient!=nil)
        
    return server, nil

}
func (s *Server) Mutation() MutationResolver {
    return &mutationResolver{s}
}

func (s *Server) Query() QueryResolver {
    return &queryResolver{s}
}

func (s *Server) Account() AccountResolver {
    return &accountResolver{s}
}

func (s *Server) Order() OrderResolver {
    return &orderResolver{s}
}
func (s *Server) Shipping() *ShippingResolver {
    return s.shippingResolver
}
type courierInfoResolver struct {
    server *Server
}
func (r *Server) CourierInfo() CourierInfoResolver {
    return &courierInfoResolver{r}
}
func (r *courierInfoResolver) Code(ctx context.Context, obj *models.CourierInfo) (string, error) {
    return obj.CourierCode, nil
}

func (r *courierInfoResolver) Name(ctx context.Context, obj *models.CourierInfo) (string, error) {
    return obj.CourierName, nil
}
type courierRateResolver struct {
    server *Server
}


func (r *courierRateResolver) CodCharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return obj.CodCharge, nil
}
func (r *courierRateResolver) FuelCharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return obj.FuelSurcharge, nil
}
func (r *courierRateResolver) ExpectedDays(ctx context.Context, obj *models.CourierRate) (int, error) {
    return (obj.ExpectedDays), nil
}
func (r *courierRateResolver) HandlingCharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return ((obj.HandlingCharge)), nil
}
func (r *courierRateResolver) TotalCharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return ((obj.TotalCharge)), nil
}


func (r *courierRateResolver) FuelSurcharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return obj.FuelSurcharge, nil
}
func (r *courierRateResolver) BaseCharge(ctx context.Context, obj *models.CourierRate) (float64, error) {
    return obj.BaseCharge, nil
}



func (r *courierInfoResolver) Description(ctx context.Context, obj *models.CourierInfo) (*string, error) {
    return &obj.Description, nil
}

func (r *courierInfoResolver) SupportedServices(ctx context.Context, obj *models.CourierInfo) ([]string, error) {
    return obj.SupportedServices, nil
}


func (s *Server) ToNewExecutableSchema() graphql.ExecutableSchema {
    return NewExecutableSchema(Config{
        Resolvers: s,
    })
}