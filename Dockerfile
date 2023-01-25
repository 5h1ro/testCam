FROM golang:1.19 AS builder

WORKDIR /go/src/app

COPY . .
RUN go mod download
RUN go mod tidy

RUN go build -o main

CMD ["./main"]