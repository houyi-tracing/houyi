#!/bin/bash

OS=$1
ARCH=$2

COMPONENT=collector
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../cmd/${COMPONENT}

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

FILTER_CONFIG=filter-config.json
cat <<EOF > ${FILTER_CONFIG}
{
  "filter-tags": []
}
EOF
mv ${FILTER_CONFIG} ${BUILD_OUT_DIR}/

RUN_COLLECTOR=run-collector.sh
cat <<EOF > ${RUN_COLLECTOR}
#!/bin/sh
./collector --sampling.max-num-child-nodes=\${MAX_NUM_CHILD_NODES} --cassandra.servers=\${CASSANDRA_SERVERS}
EOF
chmod u+x ${RUN_COLLECTOR}
mv ${RUN_COLLECTOR} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY collector /opt/ms/
COPY filter-config.json /root/
COPY ${RUN_COLLECTOR} /opt/ms/
EXPOSE 14250 14268 14269
WORKDIR=/opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_COLLECTOR}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

RUN_COLLECTOR_DOCKER=run-collector-docker.sh
cat <<EOF > ${RUN_COLLECTOR_DOCKER}
#!/bin/sh
if [[ -z \${MAX_NUM_CHILD_NODES} ]]; then
  MAX_NUM_CHILD_NODES=4
fi
if [[ -z \${CASSANDRA_SERVERS} ]]; then
  echo "at least one cassandra server is required"
  return
fi
echo "MAX_NUM_CHILD_NODES=\${MAX_NUM_CHILD_NODES}"
echo "CASSANDRA_SERVERS=\${CASSANDRA_SERVERS}"
echo "Starting collector..."
docker run -d -p 14250:14250 -p 14268:14268 -p 14269:14269 --name houyi-collector --env CASSANDRA_SERVERS=\${CASSANDRA_SERVERS} --env MAX_NUM_CHILD_NODES=\${MAX_NUM_CHILD_NODES} houyi-collector
EOF
chmod u+x ${RUN_COLLECTOR_DOCKER}
mv ${RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/

BUILD_RUN_COLLECTOR_DOCKER=build-collector-docker.sh
cat <<EOF > ${BUILD_RUN_COLLECTOR_DOCKER}
#!/bin/sh
docker build -t houyi-collector .
EOF
chmod u+x ${BUILD_RUN_COLLECTOR_DOCKER}
mv ${BUILD_RUN_COLLECTOR_DOCKER} ${BUILD_OUT_DIR}/
