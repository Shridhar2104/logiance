package main

import (
    "log"
    "net"

    "github.com/Shridhar2104/logilo/shipment/internal/config"
    "github.com/Shridhar2104/logilo/shipment/internal/service"
    pb "github.com/Shridhar2104/logilo/shipment/proto"

    "google.golang.org/grpc"
)

func main() {
    cfg := config.NewConfig()

    lis, err := net.Listen("tcp", cfg.GRPCPort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterShipmentServiceServer(s, service.NewShipmentService(cfg))

    log.Printf("Server listening on %s", cfg.GRPCPort)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}