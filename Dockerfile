# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder


# Set the working directory for the builder
WORKDIR /app

# 1. Copy go.mod and go.sum (for dependency caching)
COPY go.mod /app
COPY go.sum /app

# Download dependencies
RUN go mod download

# 2. Copy ALL application source code, including service directories
# This fixes the "package not in std" error.
COPY . /app

COPY .env /app

# 3. Build the Go application
# The binary is named 'supernova' and built in the current directory /app
RUN CGO_ENABLED=0 GOOS=linux go build -o supernova

# --- Stage 2: Final Image (Production) ---
# Use 'scratch' for the smallest, most secure final image.
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /

# Define the port (assuming your Go code listens on 8080)
EXPOSE 8080

# Copy the compiled binary from the builder stage
# Copy it to the root of the final image
COPY --from=builder /app/supernova /supernova

# Copy the environment file. Assuming you want it at the root of the container.
COPY --from=builder /app/.env /.env

# Run the application
ENTRYPOINT ["/supernova"]