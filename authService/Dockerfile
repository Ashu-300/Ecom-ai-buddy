
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

# Set Go environment variables
WORKDIR /app



# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE SERVICE_NAME with the actual name (e.g., cartService)
COPY authService/ authService/

# 3. Copy environment file 
COPY authService/.env authService/.env

# 4. Build the Go application
# REPLACE SERVICE_NAME and BINARY_NAME
RUN go build -o /app/auth ./authService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest

# Set the binary name for the entry point


# EXPOSE the specific port for this service
EXPOSE 8081

# Copy the compiled binary and the .env file
COPY --from=builder /app/auth /auth
COPY --from=builder /app/authService/.env ./.env

# Run the application
ENTRYPOINT ["/auth"]
