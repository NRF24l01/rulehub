FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN GOPROXY="https://goproxy.cn" go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /app/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
RUN chmod +x ./main

CMD ["./main"]