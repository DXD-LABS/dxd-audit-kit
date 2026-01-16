# Multi-stage Dockerfile for DXD Audit Server
# Optimization for Railway deployment

# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o dxd-audit-server ./cmd/dxd-audit-server

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/dxd-audit-server .

# Copy migrations and API spec
COPY migrations ./migrations
COPY api ./api

# Expose port (Railway will use PORT environment variable, but 8080 is our default)
EXPOSE 8080

# Run the server
CMD ["./dxd-audit-server"]
