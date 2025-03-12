# Use Ubuntu 16.04 as the base image
FROM ubuntu:16.04

# Define Go version
ARG GO_VERSION=1.23
ARG ARCH='amd64'

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    expect \
    git \
    tar \
    gcc \
    jq \
    g++ \
    gnupg2 \
    gnupg-agent \
    debsigs \
    rpm \
    build-essential \
    software-properties-common \
    python-software-properties \
    gcc-arm-linux-gnueabi \
    dpkg-sig \
    gcc-aarch64-linux-gnu

# Install Go
# It tries to get the latest patch version of Go that matches the version specified in the Makefile
# It can be simplified once we have Renovatebot support for this repo
RUN FULL_GO_VERSION=$(curl -s "https://go.dev/dl/?mode=json" | \
    jq -r --arg VERSION "${GO_VERSION}" \
    '[.[] | select(.version | test("^go\($VERSION)\\."))] | .[0].version') && \
    curl -sSL https://golang.org/dl/${FULL_GO_VERSION}.linux-${ARCH}.tar.gz -o ${FULL_GO_VERSION}.linux-${ARCH}.tar.gz && \
    tar -C /usr/local -xzf ${FULL_GO_VERSION}.linux-${ARCH}.tar.gz && \
    rm ${FULL_GO_VERSION}.linux-${ARCH}.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"

WORKDIR /app

COPY build /app/build
COPY tools /app/tools

COPY go.mod go.sum Makefile /app/

# Optional: Set Go environment flags
ENV GOFLAGS="-buildvcs=false"

# Optional: Configure git
RUN git config --global --add safe.directory /go/src/github.com/newrelic/nri-flex

RUN make deps
