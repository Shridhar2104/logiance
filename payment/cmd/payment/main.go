package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/Shridhar2104/logilo/payment"
	"github.com/Shridhar2104/logilo/payment/pb"
)

const (
	httpPort    = ":8080"
	grpcPort    = ":50051"
	dbConnString = "postgres://username:password@localhost:5432/payment_db?sslmode=disable"
)

func main() {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := payment.NewRepository(db)
	svc := payment.NewService(repo)
	httpServer := payment.NewServer(svc)
	mux := http.NewServeMux()
	httpServer.RegisterRoutes(mux)

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, payment.NewGRPCServer(svc))

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
