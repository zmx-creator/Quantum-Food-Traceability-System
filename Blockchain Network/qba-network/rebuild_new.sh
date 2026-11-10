#!/bin/bash

# Switch to the target directory
cd ~/go/src/github.com/hyperledger/fabric

# Disable Go Modules; build in GOPATH mode
echo "Disabling Go Modules, building in GOPATH mode..."
export GO111MODULE=off

# Remove all Docker images
# echo "Removing all Docker images..."
# docker rmi -f $(docker images -qa) 2>/dev/null || echo "No images to remove"

# Run make clean
echo "Running make clean..."
make clean

# Remove stale binaries
echo "Removing stale binaries..."

rm /home/test/go/src/github.com/hyperledger/fabric/.build/bin/configtxgen
rm /home/test/go/src/github.com/hyperledger/fabric/.build/docker/bin/orderer

# Build the orderer image
echo "Building the orderer image..."
make orderer-docker

# Build configtxgen
echo "Building configtxgen..."
make configtxgen

# Install the freshly built binaries
echo "Installing binaries..."
cd /home/test/fabric-samples/bin
cp /home/test/go/src/github.com/hyperledger/fabric/.build/bin/configtxgen /home/test/fabric-samples/bin
cp /home/test/go/src/github.com/hyperledger/fabric/.build/docker/bin/orderer /home/test/fabric-samples/bin

echo "All steps completed!"