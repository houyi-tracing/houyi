#!/bin/bash

ROOT_CMD='./collector'

if [[ -z $CASSANDRA_SERVERS ]]; then
  echo "\$CASSANDRA_SERVERS must be set"
  exit -1
else
  ROOT_CMD="${ROOT_CMD} --cassandra.servers=${CASSANDRA_SERVERS}"
fi

if [[ -n $LOG_LEVEL ]]; then
  ROOT_CMD="${ROOT_CMD} --log-level=${LOG_LEVEL}"
fi

if [[ -n $MAX_NUM_CHILD_NODES ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.max-num-child-nodes=${MAX_NUM_CHILD_NODES}"
fi

if [[ -n $MAX_SAMPLING_RATE ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.max-samp-prob=${MAX_SAMPLING_RATE}"
fi

if [[ -n $MIN_SAMPLING_RATE ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.min-samp-prob=${MIN_SAMPLING_RATE}"
fi

if [[ -n $AMPLIFICATION_FACTOR ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.amplification-factor=${AMPLIFICATION_FACTOR}"
fi

if [[ -n $OPERATION_DURATION ]]; then
  ROOT_CMD="${ROOT_CMD} --operation.duration=${OPERATION_DURATION}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
