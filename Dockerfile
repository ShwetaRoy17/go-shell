# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o go-shell ./app/main.go


FROM ubuntu:latest

WORKDIR /root/

COPY --from=builder /app/go-shell .

CMD ["./go-shell"]

