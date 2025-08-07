# Builds and runs the EEST (execution-spec-tests) consume engine simulator
FROM ghcr.io/astral-sh/uv:python3.12-bookworm-slim

## Default fixtures/git-ref
ARG fixtures=fixtures/blockchain_tests_engine/prague
ENV FIXTURES=${fixtures}
ARG branch=feat/working-geth-test
ENV GIT_REF=${branch} 

## Clone and install EEST
RUN apt-get update && apt-get install -y git

# Allow the user to specify a branch or commit to checkout
RUN git init execution-spec-tests && \
    cd execution-spec-tests && \
    git remote add origin https://github.com/dmitriy-b/execution-spec-tests-gnosis.git && \
    git fetch --depth 1 origin $GIT_REF && \
    git checkout FETCH_HEAD;

WORKDIR /execution-spec-tests
RUN uv sync

# Cache the fixtures. This is done to avoid re-downloading the fixtures every time
# the container starts.
# If newer version of the fixtures is needed, the image needs to be rebuilt.
# Use `--docker.nocache` flag to force rebuild.
# RUN uv run consume cache --input "$FIXTURES"
RUN git clone -b feat/geth-fixtures https://github.com/NethermindEth/gnosis-fixtures.git fixtures

## Define `consume engine` entry point using the local fixtures
ENTRYPOINT uv run consume engine -v --input "$FIXTURES"
# ENTRYPOINT uv run consume engine -v --input gnosis-fixtures/blockchain_tests_engine/prague/eip2537_bls_12_381_precompiles/bls12_g1add --sim.limit="id:tests/prague/eip2537_bls_12_381_precompiles/test_bls12_g1add.py::test_call_types[fork_Prague-blockchain_test_engine_from_state_test-inf_plus_inf-call_opcode_CALLCODE-]"