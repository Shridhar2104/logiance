# Stage 1: Build the Go app
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . ./ 

# Build the Go app from the correct location
RUN GO111MODULE=on go build -o /app/graphql/graphql ./graphql

# Stage 2: Run the app
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy the binary and the .env file
COPY --from=builder /app/graphql /app/graphql
COPY .env.development /app/.env.development

RUN chmod +x /app/graphql/graphql

EXPOSE 8084

CMD ["/app/graphql/graphql"]