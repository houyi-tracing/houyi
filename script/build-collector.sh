#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=collector
BUILD_OUT_DIR=~/houyi/${COMPONENT}/
WORK_DIR=../cmd/${COMPONENT}/

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -v ${WORK_DIR}/main.go
mv ${WORK_DIR}/main ${BUILD_OUT_DIR}
mv ${WORK_DIR}/Dockerfile ${BUILD_OUT_DIR}