# Generate the ethash verification caches.
# Use a static version because this will never need to be updated.
FROM ethereum/client-go:v1.10.20 AS geth
RUN \
 /usr/local/bin/geth makecache  1 /ethash && \
 /usr/local/bin/geth makedag    1 /ethash

# Build the simulator binary
FROM golang:1-alpine AS builder
RUN apk --no-cache add gcc musl-dev linux-headers cmake make clang build-base clang-static clang-dev

# Prepare workspace.
# Note: the build context of this simulator image is the parent directory!
ADD . /source

# Build within simulator folder
WORKDIR /source/engine
RUN go build -gcflags="all=-N -l" -o ./sim .

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ADD . /
COPY --from=geth    /ethash /ethash

EXPOSE 40000

ENTRYPOINT ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./sim", "--", "serve"]
