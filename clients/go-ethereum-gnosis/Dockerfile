ARG baseimage=ghcr.io/gnosischain/geth
ARG tag=latest

# Build the go-ethereum-gnosis binary
FROM $baseimage:$tag AS builder
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

# Install dependencies
FROM alpine:latest
RUN apk add --update bash curl jq
COPY --from=builder /usr/local/bin/geth /usr/local/bin/geth

# Generate the version.txt file.
RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev 2>/dev/null | head -1 > /version.txt

# Add the default genesis file and genesis mapper script.
ADD genesis.json /genesis.json
ADD mapper.jq /mapper.jq

# Inject the startup script.
ADD geth.sh /geth.sh
RUN chmod +x /geth.sh

# Inject the enode id retriever script.
RUN mkdir /hive-bin
ADD enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8547 8551 30303 30303/udp

ENTRYPOINT ["/geth.sh"]
