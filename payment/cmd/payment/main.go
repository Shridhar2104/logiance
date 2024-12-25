package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/razorpay/razorpay-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Shridhar2104/logilo/payment"
	"github.com/Shridhar2104/logilo/payment/pb"
)

const (
	httpPort    = ":8082"
	grpcPort    = ":50051"
	dbConnString = "postgres://payment_service_user:securepassword@payment-db:5432/payment_service_db?sslmode=disable"
)

func main() {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	repo := payment.NewRepository(db)
	svc := payment.NewService(repo)
	httpServer := payment.NewServer(svc)
	mux := http.NewServeMux()
	httpServer.RegisterRoutes(mux)

	// Add endpoint to create Razorpay order
	mux.HandleFunc("/create-order", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountID string  `json:"account_id"`
			Amount    float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		razorpayClient := razorpay.NewClient("YOUR_RAZORPAY_API_KEY", "YOUR_RAZORPAY_SECRET")

		orderData := map[string]interface{}{
			"amount":   int(req.Amount * 100), // Amount in paise
			"currency": "INR",
			"receipt":  "txn_" + req.AccountID,
		}

		order, err := razorpayClient.Order.Create(orderData, nil)
		if err != nil {
			http.Error(w, "Failed to create Razorpay order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(order)
	})

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, payment.NewGRPCServer(svc))

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	go func() {
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
		}
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	log.Println("Servers stopped")
}
