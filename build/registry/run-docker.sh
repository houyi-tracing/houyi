#!/bin/sh

ROOT_CMD='docker run -d -p 14250:14250 -p 14268:14268 -p 14269:14269 --name houyi-registry'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --env LOG_LEVEL=${LOG_LEVEL}"
fi
if [[ ${GRPC_LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --env GRPC_LISTEN_PORT=${GRPC_LISTEN_PORT}"
fi
if [[ ${REFRESH_INTERVAL} ]]; then
  ROOT_CMD="${ROOT_CMD} --env REFRESH_INTERVAL=${REFRESH_INTERVAL}"
fi
if [[ ${SCALE_FACTOR} ]]; then
  ROOT_CMD="${ROOT_CMD} --env SCALE_FACTOR=${SCALE_FACTOR}"
fi

ROOT_CMD="${ROOT_CMD} houyitracing/registry"

echo $ROOT_CMD
eval $ROOT_CMD
