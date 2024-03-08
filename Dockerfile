# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/web

EXPOSE 5555

#USER nonroot:nonroot

ENTRYPOINT ["./main"]