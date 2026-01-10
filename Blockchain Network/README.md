# Blockchain Network

## Overview

This directory contains the blockchain network configuration, consensus implementation, and smart contracts for the Quantum Blockchain-based food traceability system.

## Contents

```
Blockchain Network/
├── qba-network/          # Network configuration and scripts
├── chaincode/            # Smart contracts
├── qba/                  # QBA consensus algorithm implementation
└── README.md
```

## Subdirectories

- **qba-network/**: Network configuration files, Docker compose setup, and scripts for starting/stopping the blockchain network
- **chaincode/**: Smart contract implementations (chaincode) for the food traceability system
- **qba/**: Custom consensus algorithm implementation (QBA - Quantum Byzantine Agreement)

## Requirements

- Docker
- Docker Compose
- Hyperledger Fabric v1.4
- Go 1.15+

## Quick Start

For detailed instructions on setting up and running the blockchain network, please refer to the README files in each subdirectory:

- [qba-network/README.md](./qba-network/README.md) - Network setup and configuration
- [chaincode/README.md](./chaincode/README.md) - Smart contract information
- [qba/README.md](./qba/README.md) - Consensus algorithm details

## Network Startup

1. Place the consensus algorithm files under `/go/src/github.com/hyperledger/fabric/orderer/consensus`, ensuring the folder is named `qba`
2. Navigate to the `qba-network/` directory
3. If running for the first time, execute `./rebuild_new.sh` to compile configuration files
4. Run `./start_network.sh` to start the network

## License

This project is licensed under the Apache License 2.0.
