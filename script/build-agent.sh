#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=agent
BUILD_OUT_DIR=~/houyi/${COMPONENT}/
WORK_DIR=../cmd/${COMPONENT}/

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go
mv ${WORK_DIR}/Dockerfile ${BUILD_OUT_DIR}
