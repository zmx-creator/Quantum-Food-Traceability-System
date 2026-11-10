# Wiring a custom consensus plugin into Hyperledger Fabric 1.4.4

The QBA consensus in this work is not a standalone program: it is a **pluggable consensus
plugin compiled into the orderer binary** of a Hyperledger Fabric v1.4.4 tree. The steps below
are the integration procedure originally written for the PBFT prototype; substitute `qba` for
`pbft` everywhere (see `fabric/qba_registration.patch` for the exact current diff).

## 1. Development procedure

* Patch `configtxgen` so that it accepts the custom consensus type in `configtx.yaml`:

```go
// common/tools/configtxgen/localconfig/config.go:388
switch ord.OrdererType {
    case "qba":
}
// common/tools/configtxgen/encoder/encoder.go:38
const ConsensusTypeQBA = "qba"
// common/tools/configtxgen/encoder/encoder.go:215
switch conf.OrdererType {
    case ConsensusTypeQBA:
}
```

* Register the consensus instance:

```go
// orderer/common/server/main.go:664
consenters["qba"] = qba.New()
```

* Implement the consensus interface in `orderer/consensus/consensus.go`:

```go
// Consenter - returns the Chain that implements the block-handling callbacks.
type Consenter interface {
	HandleChain(support ConsenterSupport, metadata *cb.Metadata) (Chain, error)
}
// Chain - block handling callbacks.
type Chain interface {
	Order(env *cb.Envelope, configSeq uint64) error      // handle a normal transaction
	Configure(config *cb.Envelope, configSeq uint64) error // handle a configuration transaction
	WaitReady() error        // wait until a transaction can be processed
	Errored() <-chan struct{} // signal errors
	Start()                  // initialise chain resources
	Halt()                   // release resources
}
```

* Build the orderer image (set the `orderer/peer/tools` tag accordingly):

```
$ make orderer-docker
```

* Build `configtxgen` (output: `.build/bin/configtxgen`):

```
$ make configtxgen
```

## 2. Network topology

The topology actually used by the experiments in this repository (see
`qba-network/base/docker-compose-base.yaml`):

| Role | Domain | Host port / consensus port | Organization |
| :--: | :--: | :--: | :--: |
| Orderer | orderer0.example.com | 6050 / 6070 | OrdererOrg |
| Orderer | orderer1.example.com | 6051 / 6071 | OrdererOrg |
| Orderer | orderer2.example.com | 6052 / 6072 | OrdererOrg |
| Orderer | orderer3.example.com | 6053 / 6073 | OrdererOrg |
| Orderer | orderer4.example.com | 6054 / 6074 | OrdererOrg |
| Peer | peer0.org1.example.com | 7051 | Org1MSP |

## 3. Configuration

The plugin is driven entirely by environment variables:

* `PBFT_LISTEN_PORT`: port on which this consensus instance listens
* `PBFT_NODE_ID`: identifier of this consensus node
* `PBFT_NODE_TABLE`: semicolon-separated list of peer consensus endpoints

The variable names keep the `PBFT_` prefix for historical reasons; they configure the QBA
plugin in the current build.