# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY services/catalog/ ./services/catalog/

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o catalog ./services/catalog

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/catalog .

EXPOSE 9091

CMD ["./catalog"]