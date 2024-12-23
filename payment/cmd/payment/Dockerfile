# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /src

# Copy the entire project (since it has dependencies between services)
COPY . .

# Download dependencies
RUN go mod download

# Build the payment service
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./payment/cmd/payment

# Stage 2: Run the application
FROM gcr.io/distroless/base-debian11:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /src/payment-service .

# Copy the SQL migration
COPY payment/up.sql /app/migrations/

# Expose both HTTP and gRPC ports
EXPOSE 8080 50051

# Command to run the application
CMD ["./payment-service"]