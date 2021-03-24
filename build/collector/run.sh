#!/bin/sh

ROOT_CMD='./collector'

if [[ -z $CASSANDRA_SERVERS ]]; then
  echo "\$CASSANDRA_SERVERS must be set"
  exit 0
else
  ROOT_CMD="${ROOT_CMD} --cassandra.servers=${CASSANDRA_SERVERS}"
fi

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log.level=${LOG_LEVEL}"
fi

if [[ ${COLLECTOR_GRPC_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.grpc.port=${COLLECTOR_GRPC_PORT}"
fi

# Gossip
if [[ ${GOSSIP_REGISTRY_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.registry.addr=${GOSSIP_REGISTRY_ADDR}"
fi

# Strategy Manager
if [[ ${STRATEGY_MANAGER_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.addr=${STRATEGY_MANAGER_ADDR}"
fi
if [[ ${STRATEGY_MANAGER_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.port=${STRATEGY_MANAGER_PORT}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
