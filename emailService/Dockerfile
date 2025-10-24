
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder


WORKDIR /app



# 1. Copy go.mod and go.sum from the build context root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# REPLACE emailService with the actual name (e.g., cartService)
COPY emailService/ emailService/

# 3. Copy environment file 
COPY emailService/.env emailService/.env

# 4. Build the Go application
# REPLACE emailService and BINARY_NAME
RUN go build -o /app/email ./emailService/main.go


# --- Stage 2: Final Image (Production) ---
FROM alpine:latest

# Set the binary name for the entry point

# EXPOSE the specific port for this service
EXPOSE 8087

# Copy the compiled binary and the .env file
COPY --from=builder /app/email /email
COPY --from=builder /app/emailService/.env ./.env

# Run the application
ENTRYPOINT ["/email"]

