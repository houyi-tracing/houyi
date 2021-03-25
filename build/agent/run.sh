#!/bin/sh

ROOT_CMD="./agent"

if [[ ${COLLECTOR_SERVICE_PORT_GRPC_TRACES} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.port=${COLLECTOR_SERVICE_PORT_GRPC_TRACES}"
fi

if [[ ${STRATEGY_MANAGER_SERVICE_PORT_GRPC} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.port=${STRATEGY_MANAGER_SERVICE_PORT_GRPC}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
