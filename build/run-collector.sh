#!/bin/sh

ROOT_CMD='./collector'

if [[ -z $CASSANDRA_SERVERS ]]; then
  echo "\$CASSANDRA_SERVERS must be set"
  exit 0
else
  ROOT_CMD="${ROOT_CMD} --cassandra.servers=${CASSANDRA_SERVERS}"
fi

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log-level=${LOG_LEVEL}"
fi

if [[ ${MAX_NUM_CHILD_NODES} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.max-num-child-nodes=${MAX_NUM_CHILD_NODES}"
fi

if [[ ${MAX_SAMPLING_RATE} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.max-samp-prob=${MAX_SAMPLING_RATE}"
fi

if [[ ${MIN_SAMPLING_RATE} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.min-samp-prob=${MIN_SAMPLING_RATE}"
fi

if [[ ${AMPLIFICATION_FACTOR} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.amplification-factor=${AMPLIFICATION_FACTOR}"
fi

if [[ ${OPERATION_DURATION} ]]; then
  ROOT_CMD="${ROOT_CMD} --operation.duration=${OPERATION_DURATION}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
