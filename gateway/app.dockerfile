# Build stage
FROM golang:1.25-alpine as build

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source
COPY ./gateway/ ./gateway/

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./gateway

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/gateway .

EXPOSE 8080

CMD ["./gateway"]