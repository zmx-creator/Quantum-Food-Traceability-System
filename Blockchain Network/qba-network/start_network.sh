#!/bin/bash

# Change to the specified directory
cd ~/fabric-sample/qba-network

# Disable Go Modules, use GOPATH mode
export GO111MODULE=off

# Shut down the network
echo "Shutting down the network"
sudo ./scripts/utils.sh down

# Generate corresponding network configuration files
echo "Generating network configuration files"
./scripts/gen.sh 

# Start the network
echo "Starting the network"
./scripts/utils.sh up 

# View peer0's logs
echo "Viewing peer0's logs"
docker logs peer0.org1.example.com

# View orderer0's logs
echo "Viewing orderer0's logs"
docker logs orderer0.example.com

# Clean up old crypto-config
echo "====Cleaning up old crypto-config===="
# NOTE: MSP material is intentionally NOT copied into the web backend here.
# The Spring Boot backend runs on a separate host; after regenerating crypto-config,
# sync ./crypto-config into Backend/src/main/resources/ and restart the backend.

echo "All operations completed!"
