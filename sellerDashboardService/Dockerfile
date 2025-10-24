
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder


WORKDIR /app



# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE sellerDashboardService with the actual name (e.g., cartService)
COPY sellerDashboardService/ sellerDashboardService/

# 3. Copy environment file 
COPY sellerDashboardService/.env sellerDashboardService/.env

# 4. Build the Go application
# REPLACE sellerDashboardService and BINARY_NAME
RUN go build -o /app/sellerdashboard ./sellerDashboardService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest



# EXPOSE the specific port for this service
EXPOSE 8088

# Copy the compiled binary and the .env file
COPY --from=builder /app/sellerdashboard /sellerdashboard
COPY --from=builder /app/sellerDashboardService/.env ./.env

# Run the application
ENTRYPOINT ["/sellerdashboard"]
