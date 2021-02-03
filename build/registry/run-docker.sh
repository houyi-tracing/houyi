#!/bin/sh

ROOT_CMD='docker run -d -p 14250:14250 -p 14268:14268 -p 14269:14269 --name houyi-registry'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --env LOG_LEVEL=${LOG_LEVEL}"
fi
if [[ ${LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --env LISTEN_PORT=${LISTEN_PORT}"
fi
if [[ ${REFRESH_INTERVAL} ]]; then
  ROOT_CMD="${ROOT_CMD} --env REFRESH_INTERVAL=${REFRESH_INTERVAL}"
fi
if [[ ${RANDOM_PICK} ]]; then
  ROOT_CMD="${ROOT_CMD} --env RANDOM_PICK=${RANDOM_PICK}"
fi
if [[ ${PROB_TO_R} ]]; then
  ROOT_CMD="${ROOT_CMD} --env PROB_TO_R=${PROB_TO_R}"
fi

ROOT_CMD="${ROOT_CMD} houyitracing/registry"

echo $ROOT_CMD
eval $ROOT_CMD
