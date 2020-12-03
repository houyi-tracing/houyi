#!/bin/bash

OS=$1
ARCH=$2

if [[ -z ${OS} ]]; then
  OS=linux
fi

if [[ -z ${ARCH} ]]; then
  ARCH=amd64
fi

echo "OS=${OS}"
echo "ARCH=${ARCH}"

./build-agent.sh ${OS} ${ARCH}
./build-collector.sh ${OS} ${ARCH}
