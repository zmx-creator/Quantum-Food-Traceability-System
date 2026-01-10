#!/bin/bash

# 切换到指定目录
cd ~/go/src/github.com/hyperledger/fabric

# 禁用Go Modules，使用GOPATH模式
echo "正在禁用Go Modules，使用GOPATH模式..."
export GO111MODULE=off

# 删除所有Docker镜像
# echo "正在删除所有Docker镜像..."
# docker rmi -f $(docker images -qa) 2>/dev/null || echo "没有可删除的镜像"

# 执行make clean
echo "正在执行make clean..."
make clean

# 删除旧二进制文件
echo "正在删除旧二进制文件..."

rm /home/test/go/src/github.com/hyperledger/fabric/.build/bin/configtxgen
rm /home/test/go/src/github.com/hyperledger/fabric/.build/docker/bin/orderer

# 编译orderer镜像
echo "正在编译orderer镜像..."
make orderer-docker

# 编译configtxgen
echo "正在编译configtxgen..."
make configtxgen

# 替换二进制文件
echo "正在替换二进制文件..."
cd /home/test/fabric-samples/bin
cp /home/test/go/src/github.com/hyperledger/fabric/.build/bin/configtxgen /home/test/fabric-samples/bin
cp /home/test/go/src/github.com/hyperledger/fabric/.build/docker/bin/orderer /home/test/fabric-samples/bin

echo "所有操作已完成！"