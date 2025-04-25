# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/
COPY --from=builder /app/docs ./docs/

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 