ARG GO_VERSION=1.19

FROM golang:$GO_VERSION

WORKDIR /app

COPY build /app/build
COPY tools /app/tools

COPY go.mod go.sum Makefile /app/

RUN make deps
