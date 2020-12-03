#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=agent
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}

CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

RUN_AGENT=run-agent.sh
cat <<EOF > ${RUN_AGENT}
#!/bin/sh
nohup ./agent --reporter.grpc.host-port=${HOUYI_COLLECTOR_HOST}:14250 --collector.host=${HOUYI_COLLECTOR_HOST} &
EOF
chmod u+x ${RUN_AGENT}
mv ${RUN_AGENT} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY agent /opt/ms/
COPY run-agent.sh /opt/ms/
EXPOSE 6832 6831 5775 5778 14271
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/run-agent.sh"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_AGENT_DOCKER=run-agent-docker.sh
cat <<EOF > ${RUN_AGENT_DOCKER}
#!/bin/bash
docker run -d -p 5775:5775 -p 5778:5778 -p 6831:6831 -p 6832:6832 -p 14271:14271 --name houyi-agent --env HOUYI_COLLECTOR_HOST=${HOUYI_COLLECTOR_HOST} houyi-agent
EOF
chmod u+x ${RUN_AGENT_DOCKER}
mv ${RUN_AGENT_DOCKER} ${BUILD_OUT_DIR}/

BUILD_AGENT_DOCKER=build-agent-docker.sh
cat <<EOF > ${BUILD_AGENT_DOCKER}
#!/bin/sh
docker build -t houyi-agent .
EOF
chmod u+x ${BUILD_AGENT_DOCKER}
mv ${BUILD_AGENT_DOCKER} ${BUILD_OUT_DIR}/
