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

if [[ ${SST_MAX_CHILD_NODES} ]]; then
  ROOT_CMD="${ROOT_CMD} --sst.max.child.nodes=${SST_MAX_CHILD_NODES}"
fi
if [[ ${COLLECTOR_GRPC_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --collector.grpc.port=${COLLECTOR_GRPC_PORT}"
fi

# Gossip
if [[ ${GOSSIP_SEED_GRPC_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.seed.grpc.port=${GOSSIP_SEED_GRPC_PORT}"
fi
if [[ ${GOSSIP_REGISTRY_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.registry.addr=${GOSSIP_REGISTRY_ADDR}"
fi
if [[ ${GOSSIP_REGISTRY_GRPC_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.registry.grpc.port=${GOSSIP_REGISTRY_GRPC_PORT}"
fi
if [[ ${GOSSIP_SEED_LRU_SIZE} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.seed.lru.size=${GOSSIP_SEED_LRU_SIZE}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
