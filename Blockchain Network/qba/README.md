# Blockchain Network

## Contents
```
	.
	├── qba-network/            # network config
	│    ├── base/
	│    ├── benchmarks/
	│    ├── QBAkey/
	│    ├── scripts/
	│    ├── start_network.sh
	│    ├── rebuild_new.sh
	│    └── README.md
	├── chaincode/
	│    ├── demo/callback/
	│    └── source-app.go
	└── qba/                      # QBA consensus
   		 ├── cmd/
   		 │   └── cmd.go
   		 ├── message/
   		 ├── node/
   		 ├── server/
   		 ├── consensus.go
   		 ├── consenter.go
   		 ├── chain.go
   		 └── README.md
		 
```
## Requirements

- Fabric v1.4.4
- Go 1.15

This system uses **Docker & Docker Compose **to run the test network, and need common build tools such as make, jq, etc.

## User guide

The consensus algorithm is designed based on the QBA algorithm principle described in the paper "Protocol Description of QBA". It is implemented by referencing the original PBFT consensus algorithm example and adapting the basic elements of consensus algorithms accordingly.

### **qba-network** is the Network config folder

### **chaincode** folder contains smart contracts.

### For qba folder
**consensus.go, consenter.go**, and **chain.go** are interfaces for implementing pluggable consensus in Fabric.

1. **cmd/** contains the basic configuration definitions for the network.
2. **server/** configures the message interfaces. Based on this structure, a series of message interfaces are redesigned specifically for the QBA algorithm.
3. **node/** handles node configuration and the consensus section, which is primarily designed according to the workflow of the QBA algorithm.

### For network startup steps, please refer to the **README.md** file in the root directory.

## Consensus plugin modules

| Module | Purpose |
| --- | --- |
| `consensus.go`, `consenter.go`, `chain.go` | Fabric pluggable-consensus entry points; `Order()` forwards a request to the primary |
| `cmd/cmd.go` | reads `PBFT_NODE_ID`, `PBFT_NODE_TABLE`, `PBFT_LISTEN_PORT`, `QBA_ROUND_TIMEOUT_MS`, `QBA_MAX_RETRY` |
| `message/` | message types and the QDS construction (Toeplitz universal hash over three key lines) |
| `node/qbasend1.go`, `qbasend2.go`, `qbaforward.go`, `qbaverify.go`, `qbacommit.go`, `reply.go` | the phases of a round |
| `node/utils.go` | key-pool access and the round helper functions; `GetPrimary()` currently always returns node 0 |
| `node/round.go` | round deadline, bounded retransmission and safe failure; keeps every node on the primary key window |
| `node/faultinject.go` | runtime `tamper` fault mode for the demonstration |
| `node/trace.go` | per-round QDS step trace behind `GET /trace` |
| `server/` | the HTTP listener, including the demonstration administration endpoints |

## Administration endpoints (demonstration)

`POST /fault?mode=none|tamper`, `POST /faultpeers?ids=`, `GET /status` and `GET /trace?since=&limit=`
are served on `PBFT_LISTEN_PORT` alongside the consensus message endpoints. They
only switch the fault mode of the node and report state; the honest path returns
from both hooks immediately when the mode is `none`.

Private key material is never logged and never placed in the trace: the logs carry
key line indices and a short digest fingerprint only.
