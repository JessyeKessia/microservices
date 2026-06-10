FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

# O truque está no asterisco (go.su*). Se não houver go.sum, ele não quebra!
COPY go.mod go.su* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o order-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order-service .

EXPOSE 3000

CMD ["./order-service"]
