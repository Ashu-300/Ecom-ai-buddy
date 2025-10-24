
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

# Set Go environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux
WORKDIR /app

# Install necessary build tools for CGO on Alpine (Fixes "gcc not found")
RUN apk update && apk add --no-cache build-base

# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE SERVICE_NAME with the actual name (e.g., cartService)
COPY orderService/ orderService/

# 3. Copy environment file 
COPY orderService/.env orderService/.env

# 4. Build the Go application
# REPLACE SERVICE_NAME and BINARY_NAME
RUN go build -o /app/order ./orderService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest



# EXPOSE the specific port for this service
EXPOSE 8084

# Copy the compiled binary and the .env file
COPY --from=builder /app/order /order
COPY --from=builder /app/orderService/.env ./.env

# Run the application
ENTRYPOINT ["/order"]
