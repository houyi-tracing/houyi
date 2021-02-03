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

./agent/build.sh ${OS} ${ARCH}
./registry/build.sh ${OS} ${ARCH}
./sm/build.sh ${OS} ${ARCH}
./collector/build.sh ${OS} ${ARCH}
