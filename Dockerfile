# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

ADD . /app/

RUN go build ./cmd/web

EXPOSE 5555

USER nonroot:nonroot

ENTRYPOINT ["main.go"]