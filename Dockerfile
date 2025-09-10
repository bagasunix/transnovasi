# ---- Build stage ----
FROM golang:alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files and download deps first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app ./cmd/api

# ---- Runtime stage ----
FROM ubuntu:latest

WORKDIR /app

# Set timezone
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime \
    && echo "Asia/Jakarta" > /etc/timezone \
    && apk del tzdata

# Copy binary from builder
COPY --from=builder /app/bin/app .

# Expose port
EXPOSE 8080

# Run binary
CMD ["./app"]
