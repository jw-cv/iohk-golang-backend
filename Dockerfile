# Build Stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env.local ./

EXPOSE ${APP_PORT:-8080}

CMD ["./main"]