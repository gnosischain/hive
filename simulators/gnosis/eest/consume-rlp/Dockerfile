# Builds and runs the EEST (execution-spec-tests) consume rlp simulator
FROM python:3.10-slim

## Default fixtures
ARG fixtures=latest-stable-release
ENV INPUT=${fixtures}
ARG branch=main
ENV BRANCH=${branch}

## Install dependencies
RUN apt-get update && \
    apt-get install -y git wget tar && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

## Clone and install EEST
RUN git clone https://github.com/gnosischain/execution-spec-tests.git --branch "$BRANCH" --single-branch --depth 1
WORKDIR execution-spec-tests
RUN pip install uv && uv sync

## Define `consume rlp` entry point using the local fixtures
ENTRYPOINT uv run consume rlp -v --input "$INPUT"
