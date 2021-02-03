#!/bin/sh

ROOT_CMD="./agent"

if [[ ${COLLECTOR_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.addr=${COLLECTOR_ADDR}"
fi
if [[ ${COLLECTOR_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.port=${COLLECTOR_PORT}"
fi
if [[ ${STRATEGY_MANAGER_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.addr=${STRATEGY_MANAGER_ADDR}"
fi
if [[ ${STRATEGY_MANAGER_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.port=${STRATEGY_MANAGER_PORT}"
fi
if [[ ${GRPC_LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --grpc.listen.port=${GRPC_LISTEN_PORT}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
