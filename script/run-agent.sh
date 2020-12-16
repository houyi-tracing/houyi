#!/bin/sh

ROOT_CMD="./agent"

if [[ $COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=${HOUYI_COLLECTOR_SERVICE_HOST}:14250 --collector.host=${HOUYI_COLLECTOR_SERVICE_HOST}"
else
  echo "\$COLLECTOR_SERVICE_HOST must be set"
  exit -1
fi

echo $ROOT_CMD
eval $ROOT_CMD
