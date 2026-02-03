ARG baseimage=ghcr.io/gnosischain/reth_gnosis
ARG tag=prerelease-v1

# Build the reth-gnosis binary
FROM $baseimage:$tag AS builder
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

# Install script tools.
RUN apt-get update -y
RUN apt-get install -y bash curl jq

# Add default genesis file and genesis mapper script.
COPY genesis.json /genesis.json
COPY mapper.jq /mapper.jq

# Add the startup script.
COPY reth.sh /reth.sh
RUN chmod +x /reth.sh

# Add the enode URL retriever script.
RUN mkdir /hive-bin
COPY enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

# Create version.txt
RUN /usr/local/bin/reth --version | sed -e 's/reth \(.*\)/\1/' > /version.txt

# Export the usual networking ports to allow outside access to the node.
EXPOSE 8545 8546 30303 30303/udp

ENTRYPOINT ["/reth.sh"]