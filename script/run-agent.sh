#!/bin/sh

ROOT_CMD="./agent"

if [[ -n $COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=${COLLECTOR_SERVICE_HOST}:14250 --collector.host=${COLLECTOR_SERVICE_HOST}"
else
  echo "\$COLLECTOR_SERVICE_HOST must be set"
  exit -1
fi

echo $ROOT_CMD
eval $ROOT_CMD
