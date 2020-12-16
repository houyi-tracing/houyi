#!/bin/sh

ROOT_CMD="./agent"

if [[ $COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=http://houyi-collector:14250 --collector.host=houyi-collector"
else
  echo "\$COLLECTOR_SERVICE_HOST must be set"
  exit -1
fi

echo $ROOT_CMD
eval $ROOT_CMD
