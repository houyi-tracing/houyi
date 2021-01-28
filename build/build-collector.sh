#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=collector
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

FILTER_CONFIG=filter-config.json
cp ${FILTER_CONFIG} ${BUILD_OUT_DIR}/

RUN_COLLECTOR=run-collector.sh
chmod u+x ${RUN_COLLECTOR}
cp ${RUN_COLLECTOR} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY collector /opt/ms/
COPY filter-config.json /root/
COPY ${RUN_COLLECTOR} /opt/ms/
EXPOSE 14250 14268 14269
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_COLLECTOR}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_COLLECTOR_DOCKER=run-collector-docker.sh
chmod u+x ${RUN_COLLECTOR_DOCKER}
cp ${RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/

BUILD_RUN_COLLECTOR_DOCKER=build-collector-docker.sh
chmod u+x ${BUILD_RUN_COLLECTOR_DOCKER}
cp ${BUILD_RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/
