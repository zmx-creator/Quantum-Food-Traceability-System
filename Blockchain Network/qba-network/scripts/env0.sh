#!/bin/sh

K=5

PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations


ORDERER0NODE=orderer0.example.com:6050
ORDERER1NODE=orderer1.example.com:6051
ORDERER2NODE=orderer2.example.com:6052
ORDERER3NODE=orderer3.example.com:6053
ORDERER4NODE=orderer4.example.com:6054
ORDERER5NODE=orderer5.example.com:6055

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

    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c '{"Args":["addProInfo","FOOD001","apple","500g","2025-01-01","2025-12-31","LOT001","QS001","TestCompany","CNY 10.00","Liaoning"]}'

 
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c '{"Args":["addIngInfo","FOOD001","ING001","apple","ING002","suger"]}'


    peer chaincode invoke \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        --peerAddresses ${PEERORG1NODE} \
        -c '{"Args":["addLogInfo","FOOD001","2025-01-01 08:00","2025-01-01 18:00","Transport","Liaoning","Nanning","NanningSupermarket","2025-01-01 18:00","Plane","SF","CNY 50.00"]}'



    echo "All tests completed!"
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
    all)
        Org1
        InstallChannel
        JoinChannel
        AnchorUpdate
        InstallChainCode
        InstantiateChainCode
        #TestDemo
        ;;
esac
