ARG baseimage=ghcr.io/gnosischain/reth_gnosis
ARG tag=master

FROM $baseimage:$tag as builder
ARG GOPROXY

ENV GOPROXY=${GOPROXY}

# Install script tools.
RUN apt-get update -y
RUN apt-get install -y bash curl jq

# Add genesis mapper script.
ADD genesis.json /genesis.json
ADD mapper.jq /mapper.jq

# Add the startup script.
ADD reth.sh /reth.sh
RUN chmod +x /reth.sh

# Add the enode URL retriever script.
ADD enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

# Create version.txt
RUN /usr/local/bin/reth --version | sed -e 's/reth \(.*\)/\1/' > /version.txt

# Export the usual networking ports to allow outside access to the node.
EXPOSE 8545 8546 30303 30303/udp

ENTRYPOINT ["/reth.sh"]