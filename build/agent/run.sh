#!/bin/sh

ROOT_CMD="./agent --config.server.addr=config-server --collector.addr=collector"

if [[ ${COLLECTOR_SERVICE_PORT_GRPC_TRACES} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.port=${COLLECTOR_SERVICE_PORT_GRPC_TRACES}"
fi

if [[ ${CONFIG_SERVER_SERVICE_PORT_GRPC} ]]; then
  ROOT_CMD="${ROOT_CMD} --config.server.port=${CONFIG_SERVER_SERVICE_PORT_GRPC}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
