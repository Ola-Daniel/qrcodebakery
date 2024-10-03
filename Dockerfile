# syntax=docker/dockerfile:1
FROM golang:latest AS builder
RUN ldd --version
WORKDIR /build
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y upgrade bsdutils
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
WORKDIR /build/cmd/web
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o qrcodebakery


FROM alpine:latest
RUN lld; exit 0
WORKDIR /app
COPY --from=builder /build/cmd/web/qrcodebakery .
EXPOSE 5555
#USER nonroot:nonroot
ENTRYPOINT ["/app/qrcodebakery"]