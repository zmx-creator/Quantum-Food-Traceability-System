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


