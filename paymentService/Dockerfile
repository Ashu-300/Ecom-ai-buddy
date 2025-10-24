

# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder


WORKDIR /app



# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE SERVICE_NAME with the actual name (e.g., cartService)
COPY paymentService/ paymentService/

# 3. Copy environment file 
COPY paymentService/.env paymentService/.env

# 4. Build the Go application
# REPLACE SERVICE_NAME and BINARY_NAME
RUN go build -o /app/payment ./paymentService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest



# EXPOSE the specific port for this service
EXPOSE 8085

# Copy the compiled binary and the .env file
COPY --from=builder /app/payment /payment
COPY --from=builder /app/paymentService/.env ./.env

# Run the application
ENTRYPOINT ["/payment"]

