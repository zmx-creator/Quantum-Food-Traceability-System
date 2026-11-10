#!/bin/bash

ORG1=org1.example
ORG1USERS=(Admin User1)
VERSION=qba

# Copy the keystore key into key.pem
CPFile() {
    files=$(ls $1)
    echo ${files[0]}
    cd $1
    cp ${files[0]} ./key.pem
    cd -
}

# Do that for every Org1 user keystore
CPAllFiles() {
    PREFIX=crypto-config/peerOrganizations
    SUFFIX=msp/keystore
    for u in ${ORG1USERS[@]}; do
        CPFile ${PREFIX}/${ORG1}.com/users/${u}@${ORG1}.com/${SUFFIX}
    done
}

# Remove generated state
Clean() {
    rm -rf ./channel-artifacts
    rm -rf ./crypto-config
    rm -rf ./production
    rm -rf /tmp/crypto
}

case $1 in
    # start/stop used by the stress test
    cli)
        env IMAGETAG=${VERSION} docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec -ti cli /bin/bash
        ;;
    up)
        CPAllFiles
        env IMAGETAG=${VERSION} docker-compose -f ./docker-compose-cli.yaml up -d
        sleep 3
        docker exec cli /bin/bash -c "scripts/env.sh all"
        ;;
    down)
        docker kill $(docker ps -qa)
        echo y | docker system prune
        docker rmi $(docker images | grep 'dev-*' | awk '{print $3}')
        echo y | docker system prune
        Clean
        ;;
esac
