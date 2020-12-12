#!/bin/sh

ROOT_CMD="docker run -d -p 5775:5775 -p 5778:5778 -p 6831:6831 -p 6832:6832 -p 14271:14271 --name houyi-agent"

if [[ $COLLECTOR_SERVICE_HOST ]]; then
  ROOT_CMD="${ROOT_CMD} --env COLLECTOR_SERVICE_HOST=${COLLECTOR_SERVICE_HOST}"
else
  echo "\$COLLECTOR_SERVICE_HOST must be set"
  exit -1
fi

ROOT_CMD="$ROOT_CMD houyitracing/agent"

echo $ROOT_CMD
eval $ROOT_CMD
