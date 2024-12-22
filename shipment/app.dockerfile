# ./shipment/app.dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /shipment-service ./shipment/cmd/server

# Create final lightweight image
FROM alpine:latest

# Add ca certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /shipment-service .

# Expose gRPC port
EXPOSE 50051

# The environment variables will be passed from docker-compose
ENV DELHIVERY_API_KEY=""
ENV DELHIVERY_BASE_URL=""
ENV BLUEDART_API_KEY=""
ENV BLUEDART_BASE_URL=""

# Run the service
CMD ["./shipment-service"]