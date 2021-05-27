ARG GO_VERSION=1.15

FROM golang:$GO_VERSION

WORKDIR /app

COPY build /app/build
COPY tools /app/tools

COPY go.mod go.sum Makefile /app/

RUN make deps
