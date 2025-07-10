# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o filmhub ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/filmhub ./filmhub
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/pkg ./pkg
COPY --from=builder /app/vendor ./vendor
COPY --from=builder /app/go.mod ./go.mod
COPY --from=builder /app/go.sum ./go.sum

ENV DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_HOST=postgres \
    DB_PORT=5432 \
    DB_NAME=filmhub

EXPOSE 8080
CMD ["./filmhub"]
