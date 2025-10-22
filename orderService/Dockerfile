# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
WORKDIR /app

# 1. Copy go.mod and go.sum from the build context root (which is the main project root)
# The source path is relative to the directory where the 'docker build' command is executed.
# Since we execute from the root, 'go.mod' refers to the file at the root.
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# 2. Copy the entire service directory into the container
# Since we execute from the project root, 'authService' is a valid directory path.
# Change 'authService' to 'cartService', 'orderService', etc., for the other service Dockerfiles.
# The destination is /app/authService
COPY orderService/ orderService/

COPY .env .env


# 3. Build the Go application, specifying the path to the service's main.go
# This builds the binary for the specific service only.
RUN go build -o /app ./orderService/main.go


# --- Stage 2: Final Image ---
FROM scratch

# IMPORTANT: CHANGE THIS PORT (e.g., 8081 for auth, 8082 for cart, etc.)
EXPOSE 8084

# Copy the compiled binary
COPY --from=builder /app /app

# Run the application
ENTRYPOINT ["/app"]
