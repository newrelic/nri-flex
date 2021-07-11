ARG buildeImage=nri-flex-builder
FROM ${buildeImage} AS base

ENV CGO_ENABLED 0

COPY . .

FROM base as httpBuilder

RUN go build -ldflags '-w -extldflags "-static"' -o bin/http-server integration-test/https-server/server.go

FROM scratch AS httpServer
ADD integration-test/https-server/cabundle /cabundle
COPY --from=httpBuilder /app/bin/http-server /
