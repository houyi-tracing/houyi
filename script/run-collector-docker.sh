#!/bin/sh

ROOT_CMD='docker run -d -p 14250:14250 -p 14268:14268 -p 14269:14269 --name houyi-collector'

if [[ -z $CASSANDRA_SERVERS ]]; then
  echo "\$CASSANDRA_SERVERS must be set"
  exit -1
else
  ROOT_CMD="${ROOT_CMD} --env CASSANDRA_SERVERS=${CASSANDRA_SERVERS}"
fi

if [[ -n $LOG_LEVEL ]]; then
  ROOT_CMD="${ROOT_CMD} --env LOG_LEVEL=${LOG_LEVEL}"
fi

if [[ -n $MAX_NUM_CHILD_NODES ]]; then
  ROOT_CMD="${ROOT_CMD} --env MAX_NUM_CHILD_NODES=${MAX_NUM_CHILD_NODES}"
fi

if [[ -n $MAX_SAMPLING_RATE ]]; then
  ROOT_CMD="${ROOT_CMD} --env MAX_SAMPLING_RATE=${MAX_SAMPLING_RATE}"
fi

if [[ -n $MIN_SAMPLING_RATE ]]; then
  ROOT_CMD="${ROOT_CMD} --env MIN_SAMPLING_RATE=${MIN_SAMPLING_RATE}"
fi

if [[ -n $AMPLIFICATION_FACTOR ]]; then
  ROOT_CMD="${ROOT_CMD} --env AMPLIFICATION_FACTOR=${AMPLIFICATION_FACTOR}"
fi

if [[ -n $OPERATION_DURATION ]]; then
  ROOT_CMD="${ROOT_CMD} --env OPERATION_DURATION=${OPERATION_DURATION}"
fi

ROOT_CMD="${ROOT_CMD} houyitracing/collector"

echo $ROOT_CMD
eval $ROOT_CMD
