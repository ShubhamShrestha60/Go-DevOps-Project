# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
# Copy migrations and static assets
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
