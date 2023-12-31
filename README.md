<!--
parent:
  order: false
-->

<div align="center">
  <h1> Astra </h1>
</div>

<div align="center">
  <a href="https://codecov.io/gh/AstraProtocol/astra">
    <img alt="Code Coverage" src="https://codecov.io/gh/AstraProtocol/astra/branch/main/graph/badge.svg" />
  </a>
</div>

Astra is a scalable, high-throughput Proof-of-Stake blockchain that is fully compatible and
interoperable with Ethereum. It's built using the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) which runs on top of [Tendermint Core](https://github.com/tendermint/tendermint) consensus engine.

**Note**: Requires [Go 1.18.0+](https://golang.org/dl/)

## Installation

For prerequisites and detailed build instructions please read the Installation instructions. Once the dependencies are installed, run:

```bash
make install
```

## Integration test
    cd integration_tests
### Install dependencies
    pip3 install -r requirements.txt 
### Run test
    ./test.sh

### Run single test
    example:
        pytest -m staking -vv


Or check out the latest [release](https://github.com/AstraProtocol/astra/releases).
