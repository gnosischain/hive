# Generate the ethash verification caches.
# Use a static version because this will never need to be updated.
FROM ethereum/client-go:v1.10.20 AS geth
RUN \
 /usr/local/bin/geth makecache  1     /ethash && \
 /usr/local/bin/geth makecache  30000 /ethash && \
 /usr/local/bin/geth makedag    1     /ethash && \
 /usr/local/bin/geth makedag    30000 /ethash

# This simulation runs Engine API tests.
FROM golang:1-alpine as builder
RUN apk add --update gcc musl-dev linux-headers

# Build the simulator executable.
ADD . /source
WORKDIR /source
RUN go build -v .

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ADD . /source
WORKDIR /source
COPY --from=geth    /ethash /ethash

EXPOSE 40000

ENTRYPOINT ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./engine", "--", "serve"]
