FROM golang:1.22-alpine AS builder

RUN mkdir -p /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY config ./config

EXPOSE 9090

CMD ["./main"]