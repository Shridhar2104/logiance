// client/client.go
package client

import (
    "context"
    
    "github.com/yourusername/shipment-service/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type ShipmentClient struct {
    conn   *grpc.ClientConn
    client proto.ShipmentServiceClient
}

func NewShipmentClient(address string) (*ShipmentClient, error) {
    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }

    client := proto.NewShipmentServiceClient(conn)
    return &ShipmentClient{
        conn:   conn,
        client: client,
    }, nil
}

func (c *ShipmentClient) Close() error {
    return c.conn.Close()
}

func (c *ShipmentClient) CalculateShippingRate(ctx context.Context, req *proto.ShippingRateRequest) (*proto.ShippingRateResponse, error) {
    return c.client.CalculateShippingRate(ctx, req)
}