#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=collector
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go
cp ${WORK_DIR}/Dockerfile ${BUILD_OUT_DIR}/
cp ${WORK_DIR}/filter-config.json ${BUILD_OUT_DIR}/
