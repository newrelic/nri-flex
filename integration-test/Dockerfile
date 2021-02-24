FROM golang:1.13-alpine3.11

ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod vendor

COPY . .
