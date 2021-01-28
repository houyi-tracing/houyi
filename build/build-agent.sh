#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=agent
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}

CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

RUN_AGENT=run-agent.sh
chmod u+x ${RUN_AGENT}
cp ${RUN_AGENT} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY agent /opt/ms/
COPY ${RUN_AGENT} /opt/ms/
EXPOSE 6832 6831 5775 5778 14271
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_AGENT}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_AGENT_DOCKER=run-agent-docker.sh
chmod u+x ${RUN_AGENT_DOCKER}
cp ${RUN_AGENT_DOCKER} ${BUILD_OUT_DIR}/

BUILD_AGENT_DOCKER=build-agent-docker.sh
chmod u+x ${BUILD_AGENT_DOCKER}
cp ${BUILD_AGENT_DOCKER} ${BUILD_OUT_DIR}/

