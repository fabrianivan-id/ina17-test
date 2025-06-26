# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download

# Build user-service
FROM builder AS user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/user-service

# Build booking-service
FROM builder AS booking-service
RUN CGO_ENABLED=0 GOOS=linux go build -o booking-service ./cmd/booking-service

# Build payment-service
FROM builder AS payment-service
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/payment-service

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=user-service /app/user-service .
COPY --from=booking-service /app/booking-service .
COPY --from=payment-service /app/payment-service .