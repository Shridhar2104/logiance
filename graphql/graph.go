package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Shridhar2104/logilo/account"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
	"github.com/Shridhar2104/logilo/shopify"
	"google.golang.org/grpc"
)
type Server struct {
    accountClient  *account.Client
    shopifyClient  *shopify.Client
    shipmentClient pb.ShipmentServiceClient
}

func NewGraphQLServer(accountUrl, shopifyUrl, shipmentUrl string) (*Server, error) {
    accountClient, err := account.NewClient(accountUrl)
    if err != nil {
        return nil, err
    }

    shopifyClient, err := shopify.NewClient(shopifyUrl)    
    if err != nil {
        accountClient.Close()
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
    }
    
    // Verify client initialization
    log.Printf("Server initialized with clients - Account: %v, Shopify: %v, Shipment: %v", 
        accountClient != nil, 
        shopifyClient != nil, 
        shipmentClient != nil)
        
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

func (s *Server) Shipping() ShippingResolver {
    return &shippingResolver{s}
}

func (s *Server) ToNewExecutableSchema() graphql.ExecutableSchema {
    return NewExecutableSchema(Config{
        Resolvers: s,
    })
}