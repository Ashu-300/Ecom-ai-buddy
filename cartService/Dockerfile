
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

# Set Go environment variables

WORKDIR /app

# Install necessary build tools for CGO on Alpine (Fixes "gcc not found")


# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE SERVICE_NAME with the actual name (e.g., cartService)
COPY cartService/ cartService/

# 3. Copy environment file 
COPY cartService/.env cartService/.env

# 4. Build the Go application
# REPLACE SERVICE_NAME and BINARY_NAME
RUN go build -o /app/cart ./cartService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest


# EXPOSE the specific port for this service
EXPOSE 8082

# Copy the compiled binary and the .env file
COPY --from=builder /app/cart /cart
COPY --from=builder /app/cartService/.env ./.env

# Run the application
ENTRYPOINT ["/cart"]
