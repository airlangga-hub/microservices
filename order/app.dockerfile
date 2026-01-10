# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY services/order/ ./services/order/

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o order ./services/order

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/order .

EXPOSE 9092

CMD ["./order"]