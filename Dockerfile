FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install swag + wire CLI
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest && \
    go install github.com/google/wire/cmd/wire@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs and wire DI, then build
RUN swag init -g cmd/server/main.go -o docs --v3.1 && \
    wire gen ./cmd/server/... && \
    go build -o bin/server ./cmd/server/...

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./server"]
