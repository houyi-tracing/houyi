#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=registry
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

RUN_SHELL=run.sh
chmod u+x ${RUN_SHELL}
cp ${RUN_SHELL} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY ${COMPONENT} /opt/ms/
COPY ${RUN_SHELL} /opt/ms/
EXPOSE 22590 22600
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_SHELL}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_SHELL_DOCKER=run-docker.sh
chmod u+x ${RUN_SHELL_DOCKER}
cp ${RUN_SHELL_DOCKER} ${BUILD_OUT_DIR}/

BUILD_RUN_SHELL_DOCKER=build-docker.sh
chmod u+x ${BUILD_RUN_SHELL_DOCKER}
cp ${BUILD_RUN_SHELL_DOCKER} ${BUILD_OUT_DIR}/
