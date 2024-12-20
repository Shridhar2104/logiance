# Build stage
FROM golang:1.21-alpine AS builder

# Add git and build tools
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Build the application
RUN cd cmd/payment && CGO_ENABLED=0 GOOS=linux go build -o payment-service

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cmd/payment/payment-service .
COPY --from=builder /app/up.sql /app/migrations/

# Install postgresql-client for database operations
RUN apk add --no-cache postgresql-client

# Expose HTTP and gRPC ports
EXPOSE 8080 50051

# Command to run the service
CMD ["./payment-service"]