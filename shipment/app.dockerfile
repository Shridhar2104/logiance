# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /src

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the shipment service
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o shipment-service ./shipment/cmd/server

# Stage 2: Run the application
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /app

# Copy the binary and necessary files
COPY --from=builder /src/shipment-service .
COPY --from=builder /src/.env.development /app/.env.development
COPY --from=builder /src/shipment/internal/database/migrations /app/internal/database/migrations

EXPOSE 8080 50051

CMD ["./shipment-service"]
