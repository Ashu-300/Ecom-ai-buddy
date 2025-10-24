
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

# Set Go environment variables

WORKDIR /app



# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE productService with the actual name (e.g., cartService)
COPY productService/ productService/

# 3. Copy environment file 
COPY productService/.env productService/.env

# 4. Build the Go application
# REPLACE productService and BINARY_NAME
RUN go build -o /app/product ./productService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest



# EXPOSE the specific port for this service
EXPOSE 8083

# Copy the compiled binary and the .env file
COPY --from=builder /app/product /product
COPY --from=builder /app/productService/.env ./.env

# Run the application
ENTRYPOINT ["/product"]
