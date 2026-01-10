#!/bin/sh


PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations


ORDERER0NODE=orderer0.example.com:6050
ORDERER1NODE=orderer1.example.com:6051
ORDERER2NODE=orderer2.example.com:6052
ORDERER3NODE=orderer3.example.com:6053
ORDERER4NODE=orderer4.example.com:6054

ORDERERNODE=${ORDERER1NODE}

PEERORG1NODE=peer0.org1.example.com:7051
CHANNEL_NAME=mychannel

NAME=source-app
VERSION=1.0


Org1(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/org1.example.com/users/Admin@org1.example.com/msp
    CORE_PEER_ADDRESS=${PEERORG1NODE}
    CORE_PEER_LOCALMSPID="Org1MSP"
    echo "node now:peer0.org1.example.com"
}


InstallChannel() {
    peer channel create \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/channel.tx \
    echo "install channel"
}


JoinChannel() {
    Org1
    peer channel join -b ${CHANNEL_NAME}.block
    echo "peer0.org1.example.com join channel" 
}


AnchorUpdate() {
    Org1
    peer channel update \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/Org1MSPanchor.tx \
    echo "org1.example update anchor peer0.org1.example.com"
}


InstallChainCode() {
    Org1
    peer chaincode install \
        -n ${NAME} \
        -v ${VERSION} \
        -p github.com/chaincode/demo/
    echo "peer0.org1.example.com install chaincode - demo"
}


InstantiateChainCode() {
    peer chaincode instantiate \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -v ${VERSION} \
        -c '{"Args":["Init"]}' \
        -P "AND ('Org1MSP.peer')"
    echo "instantiate chaincode"
    sleep 10
}



TestDemo() {
    local food_id=$1  # 接收传入的食品ID，比如 FOOD001, FOOD002...

    echo "=================start to test with ID $food_id==============="
    
    # 测试写入
    echo "1. Testing addProInfo..."
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c "{\"Args\":[\"addProInfo\",\"$food_id\",\"apple\",\"500g\",\"2025-01-01\",\"2025-12-31\",\"LOT001\",\"QS001\",\"TestCompany\",\"CNY 10.00\",\"Liaoning\"]}"

    sleep 1
    
    echo "2. Testing addIngInfo..."
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c "{\"Args\":[\"addIngInfo\",\"$food_id\",\"ING001\",\"apple\",\"ING002\",\"suger\"]}"

    sleep 1
    
    echo "3. Testing addLogInfo..."
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c "{\"Args\":[\"addLogInfo\",\"$food_id\",\"2025-01-01 08:00\",\"2025-01-01 18:00\",\"Transport\",\"Liaoning\",\"Nanning\",\"NanningSupermarket\",\"2025-01-01 18:00\",\"Plane\",\"SF\",\"CNY 50.00\"]}"  

    sleep 1

    echo "All tests for ID $food_id completed!"
}

TestDemoLoop(){
    for i in $(seq 1 100)
    do
        food_id="FOOD$(printf "%03d" $i)"  # 生成 FOOD001, FOOD002, FOOD003, FOOD004
        echo "=========== Iteration $i (ID: $food_id) ==========="
        TestDemo "$food_id"              # 将当前循环的 ID 传入 TestDemo
        echo "=========== Iteration $i completed ==========="
        sleep 3
    done
    echo "TestDemo loop completed"
}

case $1 in
    installchannel)
        InstallChannel
        ;;
    joinchannel)
        JoinChannel
        ;;
    anchorupdate)
        AnchorUpdate
        ;;
    installchaincode)
        InstallChainCode
        ;;
    instantiatechaincode)
        InstantiateChainCode
        ;;
    testdemo)
        Org1
        TestDemo
        ;;
    testdemoloop)
        Org1
        TestDemoLoop
        ;;
    all)
        Org1
        InstallChannel
        JoinChannel
        AnchorUpdate
        InstallChainCode
        InstantiateChainCode
        #TestDemoLoop
        ;;
esac
