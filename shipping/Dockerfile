FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.su* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shipping-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/shipping-service .

EXPOSE 3002

CMD ["./shipping-service"]
