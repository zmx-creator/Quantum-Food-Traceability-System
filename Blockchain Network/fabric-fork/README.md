# Fabric fork integration notes (how QBA is attached to the orderer)

The QBA consensus implemented in this project is a **custom consensus plugin of Hyperledger
Fabric v1.4.4**, compiled into the orderer binary - it is not a standalone program.

## Directory mapping

| Repository location | Virtual machine location | Notes |
|---|---|---|
| `Blockchain Network/qba/` | `~/go/src/github.com/hyperledger/fabric/orderer/consensus/qba/` | the consensus implementation (chain / consenter / node / message / server); this directory is the authoritative copy |
| `fabric-fork/fabric/qba_registration.patch` | 3 uncommitted edits in the fork working tree | registers `qba` as a valid consensus type |
| `fabric-fork/consensus-integration-notes.md` | `orderer/consensus/qba/doc.md` (historical note) | integration procedure, topology, environment variables |

## What the 9-line main-tree patch does

- `orderer/common/server/main.go`: import `orderer/consensus/qba` and set `consenters["qba"] = qba.New()` inside `initializeMultichannelRegistrar`
- `common/tools/configtxgen/localconfig/config.go`: add `case "qba"` to the `switch ord.OrdererType` initialisation, otherwise `configtxgen` panics on an unknown orderer type
- `common/tools/configtxgen/encoder/encoder.go`: add the `ConsensusTypeQBA = "qba"` constant and accept it in `NewOrdererGroup`

## Reproducing the build on a fresh machine

```bash
# 1) fetch the fabric sources at the forked baseline (commit 7917a40ff)
export GOPATH=$HOME/go; export GO111MODULE=off
git clone https://github.com/hyperledger/fabric $GOPATH/src/github.com/hyperledger/fabric
git -C $GOPATH/src/github.com/hyperledger/fabric checkout 7917a40ff

# 2) drop in the consensus package
cp -r "<repo>/Blockchain Network/qba" $GOPATH/src/github.com/hyperledger/fabric/orderer/consensus/qba

# 3) apply the main-tree registration patch
git -C $GOPATH/src/github.com/hyperledger/fabric apply "<repo>/Blockchain Network/fabric-fork/fabric/qba_registration.patch"

# 4) build the orderer image and configtxgen (same as qba-network/rebuild_new.sh)
cd $GOPATH/src/github.com/hyperledger/fabric
make clean && make orderer-docker && make configtxgen
cp .build/bin/configtxgen      ~/fabric-samples/bin/
cp .build/docker/bin/orderer   ~/fabric-samples/bin/
docker tag hyperledger/fabric-orderer:amd64-1.4.4-snapshot-7917a40ff hyperledger/fabric-orderer:latest

# 5) configtx.yaml must declare Orderer.OrdererType: qba (already set in qba-network/configtx.yaml)
```

## Verified toolchain (2026-10-26)

Go 1.18.4 / Docker 20.10.17 / Fabric 1.4.4 (commit 7917a40ff) / cryptogen + configtxgen 1.4.4

After changing anything under `Blockchain Network/qba/`, the image must be rebuilt (step 4)
before the running orderers pick the change up; a rebuild is not needed for comment-only edits.