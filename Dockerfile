# --- Stage 1: Builder ---
# Use a robust image for compiling Go applications
FROM golang:1.25-alpine AS builder

# Set necessary environment variables for CGO (required by Gin on Alpine)
ENV CGO_ENABLED=1
ENV GOOS=linux

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to fetch dependencies
COPY go.mod go.sum ./

# Download dependencies (this layer is cached if dependencies haven't changed)
RUN go mod download

# Copy the rest of the application source code
# Assuming all service directories (authService, cartService, etc.) are at the root level.
COPY . .

COPY .env .env

# Build the Go application, outputting a static binary named 'monolith'
# This command targets the root main.go file.
RUN go build -o /monolith ./main.go


# --- Stage 2: Final Image ---
# Use a minimal scratch image for the final deployment to reduce size and attack surface
FROM scratch

# Set the port the application runs on (as defined in your monolith main.go: 8080)
EXPOSE 8080

# Copy the compiled binary from the builder stage
COPY --from=builder /monolith /monolith

# Run the application
ENTRYPOINT ["/monolith"]
