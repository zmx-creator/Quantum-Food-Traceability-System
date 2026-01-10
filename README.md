# Quantum-blockchain-based food traceability system

## Overview
This project aims to demonstrate the quantum blockchain consensus algorithm through a food traceability system. It consists of three main components:
    

The Frontend component provides visual interfaces for querying and submitting data. The Backend component wraps the Hyperledger Fabric SDK to interface with the blockchain network, forwarding proposals from the front‑end to the blockchain. The blockchain network processes the proposals and returns results to the Backend component.

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

### Step for install and demo

1. Install Fabric 1.4 – follow the official [Hyperledger Fabric installation guide](https://github.com/hyperledger/fabric) steps.
2. Place the consensus algorithm files under **/go/src/github.com/hyperledger/fabric/orderer/consensus**, ensuring the folder is named **qba**. Place the **Frontend** component and **Backend** component both under **qba-network/** folder.
3. If running for the first time, execute **./rebuild_new.sh** inside the **qba-network/** folder to compile configuration files.
4. Run the script **./start_network.sh** in the **qba-network/** folder to start the network.
5. After the network is run, start the **Backend** component and the **Frontend** component in order. Refer to the respective **README.md** files for startup instructions.


## License
This project is licensed under the Apache License 2.0. For details, please see the LICENSE file.

Note: This project uses Hyperledger Fabric v1.4, which itself is licensed under the Apache License 2.0. Please ensure compliance with all relevant open‑source license terms.




