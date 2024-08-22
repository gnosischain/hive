## How tu run manually, using IDE
1. Start hive proxy
```bash
./hive --sim gnosis/go-bdd --client nethermind-gnosis --dev.addr="0.0.0.0:3000" --docker.output --dev
```
2. Build nethermind custom client
2.1. Copy `genesis.json` from tests to `clients/nethermind-gnosis/genesis.json`
```
cd clients/nethermind-gnosis
docker build -t nethermind:custom .
```
3. Start nethermind custom client

```bash
docker run -e HIVE_NETWORK_ID=10203 -e HIVE_FORK_HOMESTEAD=0 -e HIVE_FORK_BYZANTIUM=0 -e HIVE_FORK_DAO_BLOCK=0 -e HIVE_FORK_CONSTANTINOPLE=0 -e HIVE_FORK_ISTANBUL=0 -e HIVE_FORK_BERLIN=0 -e HIVE_FORK_LONDON=0 -e HIVE_TERMINAL_TOTAL_DIFFICULTY=100 -e HIVE_SHANGHAI_TIMESTAMP=1724132336 -e HIVE_CANCUN_TIMESTAMP=1724142336  -p 8545:8545 -p 8551:8551 -it nethermind:custom bash
```
4. Run tests from IDE

## How to run manually, using CLI
```
./hive --sim my-simulation --client nethermind-gnosis --dev.addr="0.0.0.0:3000" --docker.output
```
