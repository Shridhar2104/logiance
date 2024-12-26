// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/database"
	"github.com/Shridhar2104/logilo/shipment/internal/database/migrate"
	"github.com/Shridhar2104/logilo/shipment/internal/service"
	pb "github.com/Shridhar2104/logilo/shipment/proto/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
    // Configure logger to include file and line number
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Print current working directory for debugging
    pwd, _ := os.Getwd()
    log.Printf("Current working directory: %s", pwd)

    // Load environment variables
    envPath := "/app/.env.development"
    log.Printf("Loading environment from: %s", envPath)
    if err := godotenv.Load(envPath); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Load and validate configuration
    cfg := config.NewConfig()
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }

    // Initialize database connection
    db, err := initDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Get the underlying *sql.DB to run migrations and configure pool
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("Failed to get SQL DB: %v", err)
    }
    defer sqlDB.Close()

    // Configure connection pool
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // Run database migrations
    migrationsPath := filepath.Join("internal", "database", "migrations")
    log.Printf("Running migrations from: %s", migrationsPath)
    if err := migrate.RunMigrations(sqlDB, migrationsPath); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize shipment database instance
    shipmentDB := database.NewShipmentDB(db)

    // Initialize shipment service
    shipmentService, err := service.NewShipmentService(cfg, shipmentDB)
    if err != nil {
        log.Fatalf("Failed to create shipment service: %v", err)
    }

    // Create gRPC listener
    lis, err := net.Listen("tcp", cfg.GRPCPort)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    log.Printf("Server will listen on %s", cfg.GRPCPort)

    // Create gRPC server with interceptors
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(LoggingInterceptor),
    )

    // Register services
    pb.RegisterShipmentServiceServer(grpcServer, shipmentService)
    
    // Register reflection service (useful for grpcurl and other tools)
    reflection.Register(grpcServer)

    // Create channel for shutdown signals
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    // Start gRPC server in a goroutine
    go func() {
        log.Printf("Starting gRPC server on %s", cfg.GRPCPort)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()

    // Wait for shutdown signal
    sig := <-shutdown
    log.Printf("Received shutdown signal: %v", sig)

    // Initiate graceful shutdown
    log.Println("Initiating graceful shutdown...")
    grpcServer.GracefulStop()
    log.Println("Server stopped gracefully")
}

// LoggingInterceptor provides request logging and basic metrics for gRPC calls
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    startTime := time.Now()

    // Log incoming request with method name
    log.Printf("gRPC Request - Method: %s", info.FullMethod)

    // Add basic request details if available
    if stringer, ok := req.(fmt.Stringer); ok {
        log.Printf("Request details: %s", stringer.String())
    }

    // Process request
    resp, err := handler(ctx, req)

    // Calculate duration
    duration := time.Since(startTime)

    // Log response with timing and error details
    if err != nil {
        log.Printf("gRPC Error - Method: %s, Duration: %v, Error: %v",
            info.FullMethod, duration, err)
    } else {
        log.Printf("gRPC Success - Method: %s, Duration: %v",
            info.FullMethod, duration)
    }

    return resp, err
}
func initDB() (*gorm.DB, error) {
    // Get database connection string from environment
    dsn := os.Getenv("DATABASE_SHIPMENT_URL")
    if dsn == "" {
        // Use a consistent format with other services for the default connection string
        defaultDsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
            "shipment_user_dev",      // username
            "shipment_password_dev",   // password
            "localhost",            // host
            5432,                     // port
            "shipment_db_dev",        // database name
        )
        dsn = defaultDsn
    }

    // Configure GORM logger
    gormLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    // Get max connections from environment or use default
    maxConns := os.Getenv("DB_MAX_CONNECTIONS")
    maxConnections := 100 // default value
    if maxConns != "" {
        if val, err := strconv.Atoi(maxConns); err == nil {
            maxConnections = val
        }
    }

    // Open database connection with GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxOpenConns(maxConnections)
    sqlDB.SetMaxIdleConns(maxConnections / 4) // 25% of max connections
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}