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

cd ./agent
./build.sh ${OS} ${ARCH}

cd ../registry
./build.sh ${OS} ${ARCH}

cd ../sm
./build.sh ${OS} ${ARCH}

cd ../collector
./build.sh ${OS} ${ARCH}
