#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=collector
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

RUN_COLLECTOR=run.sh
chmod u+x ${RUN_COLLECTOR}
cp ${RUN_COLLECTOR} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY collector /opt/ms/
COPY ${RUN_COLLECTOR} /opt/ms/
EXPOSE 22590 22650 14580
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_COLLECTOR}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_COLLECTOR_DOCKER=run-docker.sh
chmod u+x ${RUN_COLLECTOR_DOCKER}
cp ${RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/

BUILD_RUN_COLLECTOR_DOCKER=build-docker.sh
chmod u+x ${BUILD_RUN_COLLECTOR_DOCKER}
cp ${BUILD_RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/

cp istio.yaml ${BUILD_OUT_DIR}/
cp kube.yaml ${BUILD_OUT_DIR}/
