# Start from the official Golang image for build stage
FROM golang:1.23 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Use a minimal image for the final stage
FROM alpine:latest
WORKDIR /root/

# Copy the built binary from builder
COPY --from=builder /app/server .

# Expose the port the app runs on
EXPOSE 8080

# Set environment variable for Gin release mode
ENV GIN_MODE=release

# Run the binary
CMD ["./server"] 