FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o dxd-audit-server ./cmd/dxd-audit-server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/dxd-audit-server .
CMD ["./dxd-audit-server"]
