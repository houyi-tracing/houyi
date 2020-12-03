#!/bin/bash

git pull

OS=$1
ARCH=$2

./build-agent.sh ${OS} ${ARCH}
./build-collector.sh ${OS} ${ARCH}
