// main.go
package main

import (
    "context"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Shridhar2104/logilo/shipment/internal/config"
    "github.com/Shridhar2104/logilo/shipment/internal/service"
    pb "github.com/Shridhar2104/logilo/shipment/proto"

    "google.golang.org/grpc"
)

func main() {
    // Initialize logger
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // Load configuration
    cfg := config.NewConfig()
      // Validate configuration
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    // Create listener
    lis, err := net.Listen("tcp", cfg.GRPCPort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    // Initialize shipment service
    shipmentService, err := service.NewShipmentService(cfg)
    if err != nil {
        log.Fatalf("failed to initialize shipment service: %v", err)
    }
    
    // Create gRPC server
    s := grpc.NewServer(
        grpc.UnaryInterceptor(LoggingInterceptor),
    )
    
    // Register services
    pb.RegisterShipmentServiceServer(s, shipmentService)
    
    // Create channel for shutdown signals
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
    
    // Start server in a goroutine
    go func() {
        log.Printf("Server listening on %s", cfg.GRPCPort)
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()
    
    // Wait for shutdown signal
    <-shutdown
    log.Println("Shutting down server...")
    
    // Graceful shutdown
    s.GracefulStop()
    log.Println("Server stopped")
}

// LoggingInterceptor provides request logging and basic metrics
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    
    // Log request
    log.Printf("Request - Method: %s", info.FullMethod)
    
    // Process request
    resp, err := handler(ctx, req)
    
    // Log response
    duration := time.Since(start)
    if err != nil {
        log.Printf("Response - Method: %s, Duration: %v, Error: %v", 
            info.FullMethod, duration, err)
    } else {
        log.Printf("Response - Method: %s, Duration: %v", 
            info.FullMethod, duration)
    }
    
    return resp, err
}