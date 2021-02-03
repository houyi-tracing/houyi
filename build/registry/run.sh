#!/bin/sh

ROOT_CMD='./registry'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log-level=${LOG_LEVEL}"
fi

if [[ ${LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --listen.port=${LISTEN_PORT}"
fi
if [[ ${REFRESH_INTERVAL} ]]; then
  ROOT_CMD="${ROOT_CMD} --refresh.interval=${REFRESH_INTERVAL}"
fi
if [[ ${RANDOM_PICK} ]]; then
  ROOT_CMD="${ROOT_CMD} --random.pick=${RANDOM_PICK}"
fi
if [[ ${PROB_TO_R} ]]; then
  ROOT_CMD="${ROOT_CMD} --prob.to.r=${PROB_TO_R}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
