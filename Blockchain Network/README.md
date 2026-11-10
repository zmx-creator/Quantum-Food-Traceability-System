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

## Network Topology

The demonstration network is a single-organisation Fabric v1.4 deployment in which the
five orderers are the five QBA consensus nodes:

```
N = 5 QBA consensus nodes  (orderer0 .. orderer4, custom orderer image)
    + 1 peer (peer0.org1) running the source-app chaincode
    + 1 CLI container
```

| Consensus node | Container            | Client gRPC | QBA HTTP (node-to-node) | `PBFT_NODE_ID` |
| -------------- | -------------------- | ----------- | ----------------------- | -------------- |
| Node 0         | orderer0.example.com | 6050        | 6070                    | 0              |
| Node 1         | orderer1.example.com | 6051        | 6071                    | 1              |
| Node 2         | orderer2.example.com | 6052        | 6072                    | 2              |
| Node 3         | orderer3.example.com | 6053        | 6073                    | 3              |
| Node 4         | orderer4.example.com | 6054        | 6074                    | 4              |

Role assignment inside a round:

- Node 0 is the primary for every round (`GetPrimary()` in `qba/node/utils.go` returns 0)
  and is the node that finally orders the batch; the first signer of a round is node 1
  (`qba/node/qbasend1.go`). Forwarder and verifier roles rotate within the round via
  `FindTargetID`, so the role sequence for a given configuration is deterministic.
- Fault tolerance follows `FaultNum = (len(table) - 1) / 2`, i.e. **f = 2 for N = 5**
  (`qba/cmd/cmd.go`). QBA signs every broadcast with QDS, so it belongs to the
  authenticated Byzantine agreement model where the bound is `n >= 2f + 1` rather
  than the unsigned `n >= 3f + 1`. A round tolerates up to f malicious nodes; with
  more than f the round must fail rather than commit.

Shared material:

- All five orderers mount the same QDS key directory `qba-network/QBAKey/` at `/QBAKey`.
  The key files are **not** tracked in this repository. The material used in the reported
  experiment is the output of the deployed quantum links (quantum key distribution followed
  by error correction and privacy amplification); a reproduced deployment has to supply its
  own key material in the same layout: three role files of 64-bit key lines plus the file of
  irreducible-polynomial coefficients.
- TLS is disabled on both the Fabric transport (`grpc://`) and the QBA HTTP layer. This is
  a demonstration network, not a production deployment.
### Fault model and commit rule

| | statement |
| --- | --- |
| Bound | QDS authenticates every broadcast, so the bound is `n >= 2f + 1`: with `N = 5` the network tolerates `f = 2` malicious nodes. |
| Commit rule | The primary writes a batch once **`n - f` = 3 of the 5 reports name that batch** (`Buffer.CountAgreeing` and `handleExecute` in `qba/node/qbacommit.go`). A report counts when it is a countable verdict and its decision digest equals the proposed one. |
| Verdicts | A QDS instance that fails verification is kept out of the judgment list instead of ending the round; a tie is the abstention symbol and is not counted either. Both follow section I.B step 4 and section III of the protocol description. |
| Outcome | With 1 or 2 malicious nodes the honest reports still form a quorum and the block is written; with 3 of them the honest side holds only 2 reports, no quorum is ever reached, and after `QBA_MAX_RETRY` attempts the batch is discarded as untrustworthy and no block is written. |

**Scope.** The bound is for agreement; liveness is provided by the bounded policy
described below (a fixed number of retransmissions, then the batch is discarded).
The injection targets the deputy nodes, since `GetPrimary()` designates node 0 as
the commander of every round and the scenarios demonstrated are the loyal-commander
cases of section III.B of the protocol description. Arm nodes 1 to 4.

## Subdirectories

- **qba-network/**: Network configuration files, Docker compose setup, and scripts for starting/stopping the blockchain network
- **chaincode/**: Smart contract implementations (chaincode) for the food traceability system
- **qba/**: Custom consensus algorithm implementation (QBA - Quantum Byzantine Agreement)

## Requirements

- Docker
- Docker Compose
- Hyperledger Fabric v1.4
- Go 1.15+ (verified with 1.18.4)

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

`scripts/gen.sh` generates `crypto-config/`, `channel-artifacts/` and `production/`;
`scripts/utils.sh down` tears the network down and removes those three directories, so
both the MSP material and the ledger state are destroyed. Re-generate and restart before
running another experiment.

Note on ledger persistence: `base/docker-compose-base.yaml` mounts `../production/orderer`
at `/var/hyperledger/production/orderer{N}`, while the orderer actually writes its ledger to
`/var/hyperledger/production/orderer`. Each orderer therefore keeps its ledger inside the
container filesystem, and it is lost when the container is removed. Each orderer keeps an
independent ledger, which is what makes an inter-node consistency check meaningful.

## Malicious-node demonstration

The consensus plugin exposes three administration endpoints on each node HTTP
listener (`6070..6074` on the network host, already published by the compose
file):

| Endpoint | Effect |
| --- | --- |
| `POST /fault?mode=none\|tamper` | arm or disarm this node at runtime |
| `GET /status` | mode, injection counter, round state, key window, timeout, safe-failure and consistency-violation counters |
| `GET /trace?since=&limit=` | the QDS steps of a round: phase, signer/forwarder/verifier, QDS number, digest fingerprint, outcome |

Fault modes:

- `tamper` is the malicious behaviour of the five party case in the protocol
  description: a deviating deputy whose commander is loyal.
  - as a **relay** it forwards the message and the signature it received unchanged:
    a QDS signature cannot be forged by a party outside the instance, so there is
    nothing to gain there, and the honest tables stay clean.
  - as a **signer** it signs a message other than the proposed one, and offers it
    only in the instances where the relay is itself deviating - two of the three
    parties of that instance are then deviating, which is the collusion the
    protocol analyses. In front of an honest relay it keeps to the proposed value,
    which is what the round consistency check below exists to enforce.
  - it also reports a decision for a message other than the one proposed, and the
    commit quorum simply never counts that report.
- The commander (node 0) cannot be armed: the injected strategy has no dishonest
  commander behaviour, and the backend rejects an attempt to arm it.
- Both switch at runtime. No container is restarted, so MSP material and ledger
  state survive and the demonstration can be repeated with one click.

### Consistency check between the QDS instances of a round

Protocol description of QBA, section III: in broadcast phase II every deputy
re-signs what it received, and each honest deputy compares the message it is
being shown now with the one the same deputy showed it earlier in the round. A
mismatch means the deputy is sending inconsistent messages, and it is asked to
send again.

The implementation keeps, per round, the value every signer has attested to this
node (`qba/node/consistency.go`, called from the verifier in
`qba/node/qbaverify.go` right after a successful QDS verification):

- first message of a signer: remembered, accepted;
- identical repetition: accepted, this is what an honest run always produces;
- contradiction: counted as a violation, answered to the forwarder as a failed
  verdict so the inconsistent value cannot join the judgment list, and answered
  with a resend request (`inconsistent message, resend`), capped at
  `QBA_CONSISTENCY_RETRY` requests per signer per round.

The honest path never fires, so the check costs one digest comparison per verified
message and changes nothing. It covers precisely the case QDS cannot: a signer
acting within its own key material, whose signatures are valid by construction.

Round liveness (`qba/node/round.go`): every round has a deadline
(`QBA_ROUND_TIMEOUT_MS`, default 5000 ms) and at most `QBA_MAX_RETRY` (default 5)
attempts. A retry retransmits the byte-identical signed batch and never signs
again, so no QDS key line is used for two different messages. A round that never
reaches the quorum ends in a safe failure: nothing is written, the batch is
discarded rather than re-queued, the primary broadcasts the abandonment so every
node advances its key window together, and the next client request opens a fresh
round on fresh key lines. The demonstration network sets 2000 ms and 5 attempts
in `base/peer-base.yaml`, so a round without quorum resolves in about eleven
seconds and the client is told the transaction was discarded as untrustworthy.

`scripts/byzantine_acceptance.sh` is a black-box regression gate over the whole
flow, 27 checks: an honest round commits and raises no consistency alarm; one and
two malicious nodes are outvoted and their block is still written, and a lone
deviating deputy is shown to place no message in any honest judgment list while its
dishonest report adds nothing to the quorum; three malicious nodes produce no block,
the batch is discarded after the five attempts and never reaches the ledger; the
network then recovers on the next honest round with all key windows aligned. The last
seven checks are `scripts/byzantine_cells.py` reading the recorded trace back from the
five nodes and comparing it cell by cell with the broadcast tables the protocol
states. Run it on the network host with `bash scripts/byzantine_acceptance.sh`.

The web page for this lives at `/Home/Byzantine` (see the frontend README).

## License

This project is licensed under the Apache License 2.0.
