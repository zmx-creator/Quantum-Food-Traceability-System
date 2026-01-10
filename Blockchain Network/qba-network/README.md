# QBA Network Configuration

## Overview

This directory contains the network configuration files and scripts for deploying and managing the Hyperledger Fabric blockchain network with the QBA consensus algorithm.

## Contents

```
qba-network/
├── base/                  # Base Docker configurations
├── benchmarks/            # Performance benchmarking tools
├── channel-artifacts/     # Generated channel configuration artifacts
├── crypto-config/         # Cryptographic material for network entities
├── QBAKey/               # QBA consensus algorithm keys
├── scripts/              # Utility scripts for network management
├── production/           # Production environment configurations
├── configtx.yaml         # Channel configuration
├── crypto-config.yaml    # Crypto material generation configuration
├── docker-compose-cli.yaml    # Docker Compose configuration
├── rebuild_new.sh        # Script to rebuild network artifacts
├── start_network.sh      # Script to start the network
├── start_springSDK.sh    # Script to start the Spring SDK backend
├── start_TrackWeb.sh     # Script to start the frontend web application
└── fabric_network_test.sh # Network testing script
```

## Requirements

- Docker
- Docker Compose
- Hyperledger Fabric v1.4

## Network Configuration

### Crypto Material

The `crypto-config.yaml` file defines the network topology for generating cryptographic material:

```bash
# Generate crypto material
cryptogen generate --config=crypto-config.yaml
```

### Channel Configuration

The `configtx.yaml` file defines the channel configuration, including:
- Channel profile
- Orderer configuration
- Organization configuration
- Anchor peers

## Scripts

### rebuild_new.sh

Rebuilds all network artifacts including channel configuration and crypto material. Run this when:
- First time setup
- Configuration changes have been made
- Crypto material needs to be regenerated

```bash
./rebuild_new.sh
```

### start_network.sh

Starts the blockchain network using Docker Compose:

```bash
./start_network.sh
```

This script:
1. Starts the orderer node
2. Starts peer nodes
3. Creates the channel
4. Joins peers to the channel
5. Installs and instantiates chaincode

### stop_network.sh

Stops and cleans up the network (if available):

```bash
./stop_network.sh
```

## Docker Compose Configuration

The `docker-compose-cli.yaml` file defines the network services including:
- Orderer
- Peer nodes
- CLI container

## Channel Artifacts

After running `rebuild_new.sh`, the following artifacts are generated in `channel-artifacts/`:
- Channel configuration block
- Anchor peer updates
- Genesis block (for ordering service)

## Troubleshooting

1. **Network fails to start**: Ensure Docker and Docker Compose are running
2. **Port conflicts**: Check if ports 7050-7055 are available
3. **Crypto material issues**: Run `rebuild_new.sh` to regenerate

## Integration with Backend and Frontend

After the network is started:
1. Run `start_springSDK.sh` to start the backend service
2. Run `start_TrackWeb.sh` to start the frontend application
3. Access the application at the configured endpoint

## License

This project is licensed under the Apache License 2.0.
