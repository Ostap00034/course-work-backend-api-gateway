FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o svc ./cmd/api/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/svc .
COPY .env ./
EXPOSE 8080
ENTRYPOINT ["./svc"]