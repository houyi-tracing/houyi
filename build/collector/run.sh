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

if [[ ${COLLECTOR_SERVICE_PORT_GRPC_TRACES} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.grpc.port=${COLLECTOR_SERVICE_PORT_GRPC_TRACES}"
fi

# Gossip
if [[ ${REGISTRY_SERVICE_PORT_GRPC_GOSSIP} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.registry.grpc.port=${REGISTRY_SERVICE_PORT_GRPC_GOSSIP}"
fi

# Strategy Manager
if [[ ${STRATEGY_MANAGER_SERVICE_PORT_GRPC} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.port=${STRATEGY_MANAGER_SERVICE_PORT_GRPC}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
