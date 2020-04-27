FROM golang:1.13-alpine3.11

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

ENV CGO_ENABLED 0
