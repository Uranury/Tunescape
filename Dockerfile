# ---------- Builder ----------
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/

COPY . .

ENV CGO_ENABLED=0
RUN go build -mod=vendor -o app ./cmd/api

# ---------- Runtime ----------
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080
CMD ["./app"]