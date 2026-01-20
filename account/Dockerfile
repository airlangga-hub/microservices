# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source
COPY . ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o account ./

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/account .

EXPOSE 9090

CMD ["./account"]