#!/bin/sh

ROOT_CMD="./agent"

if [[ $COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=${COLLECTOR_SERVICE_HOST}:${COLLECTOR_SERVICE_PORT_GRPC} --collector.host=${COLLECTOR_SERVICE_HOST}"
else
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=houyi-collector:14250 --collector.host=houyi-collector"
fi

echo $ROOT_CMD
eval $ROOT_CMD
