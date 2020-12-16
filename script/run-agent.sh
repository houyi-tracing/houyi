#!/bin/sh

ROOT_CMD="./agent"

if [[ $HOUYI_COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=${HOUYI_COLLECTOR_SERVICE_HOST}:14250 --collector.host=${HOUYI_COLLECTOR_SERVICE_HOST}"
else
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=houyi-collector:14250 --collector.host=houyi-collector"
fi

echo $ROOT_CMD
eval $ROOT_CMD
