# Quantum-blockchain-based food traceability system

## Overview
This project demonstrates the quantum blockchain consensus algorithm (QBA) through a food traceability system. It consists of three main components:

- **Frontend** - Vue 3 interfaces for querying and submitting data, including a *malicious-node demonstration* page (`/Home/Byzantine`) that shows a QDS-signed consensus round step by step, lets you make one to three of the five consensus nodes act malicious, shows the judgment list each node built (the m1, m2, m41 of the protocol), the consistency check that catches a signer sending different messages to different peers, the commit quorum count and whether the network agreed.
- **Backend** - a Spring Boot application wrapping the Hyperledger Fabric SDK. It forwards front-end proposals to the blockchain and exposes the administration API used by the demonstration page (`/api/fabric/demo/*`).
- **Blockchain Network** - a five-orderer Fabric v1.4 deployment in which the QBA consensus plugin is compiled into the orderer, plus the food traceability chaincode.

The blockchain network processes proposals and returns results to the Backend component.

## Contents
```
.                          
├── Frontend/                    # front end 
│	   ├── public/     
│      ├── src/       
│      ├── package.json            
│      └── README.md
│
├── Backend/                     # back end 
│	   ├── src/         
│      └── README.md
│
├── Blockchain Network/
│      ├── qba-network/          # network config
│      ├── chaincode/
│      ├── qba/                  # QBA consensus
│      └── README.md
│
└── README.md	
	
```                            

## User Guides and Executables

You can find separate User Guides in the individual README.md files for each component. Please refer to those README.md files for deployment and execution instructions.

The **chaincode/** folder contains smart contracts.

The **qba-network/** folder contains network configuration files.

### Ports

| Purpose | Port |
| --- | --- |
| Backend REST API | 8080 |
| Frontend dev server | 8081 |
| Orderer client gRPC (5 nodes) | 6050-6054 |
| QBA node-to-node HTTP and administration endpoints (5 nodes) | 6070-6074 |
| Peer | 7051 |

### Step for install and demo

1. Install Fabric 1.4 – follow the official [Hyperledger Fabric installation guide](https://github.com/hyperledger/fabric) steps.
2. Place the consensus algorithm files under **/go/src/github.com/hyperledger/fabric/orderer/consensus**, ensuring the folder is named **qba**. The **Frontend** and **Backend** components stay at the repository root; they do not need to be copied into **qba-network/**.
3. If running for the first time, execute **./rebuild_new.sh** inside the **qba-network/** folder to compile configuration files.
4. Run the script **./start_network.sh** in the **qba-network/** folder to start the network.
5. After the network is run, start the **Backend** component and the **Frontend** component in order. Refer to the respective **README.md** files for startup instructions.

The network generates its own MSP material, channel artifacts and ledger on first
start (`scripts/gen.sh`); `crypto-config/`, `channel-artifacts/`, `production/` and
the QDS key files under `qba-network/QBAKey/` are not part of this repository.

### Malicious-node demonstration

Open `/Home/Byzantine` once the network, the Backend and the Frontend are running.
Fault modes are set at runtime, so no container restarts and the ledger survives
between rounds. See
[Blockchain Network/README.md](./Blockchain Network/README.md#byzantine-demonstration)
for the fault modes, the round timeout and retry settings, and the commit quorum
(`n - f` of `n`, with `f = 2` for five nodes).


## License
This project is licensed under the Apache License 2.0. For details, please see the LICENSE file.




