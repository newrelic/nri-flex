ARG GO_VERSION=1.20

FROM golang:$GO_VERSION

WORKDIR /app

COPY build /app/build
COPY tools /app/tools

COPY go.mod go.sum Makefile /app/

ENV GOFLAGS="-buildvcs=false"
RUN git config --global --add safe.directory /go/src/github.com/newrelic/nri-flex

RUN make deps
