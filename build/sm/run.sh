#!/bin/sh

ROOT_CMD='./sm'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log.level=${LOG_LEVEL}"
fi
if [[ ${GRPC_LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --grpc.listen.port=${GRPC_LISTEN_PORT}"
fi
if [[ ${REFRESH_INTERVAL} ]]; then
  ROOT_CMD="${ROOT_CMD} --refresh.interval=${REFRESH_INTERVAL}"
fi
if [[ ${SCALE_FACTOR} ]]; then
  ROOT_CMD="${ROOT_CMD} --scale.factor=${SCALE_FACTOR}"
fi

# Gossip
if [[ ${GOSSIP_SEED_GRPC_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.seed.grpc.port=${GOSSIP_SEED_GRPC_PORT}"
fi
if [[ ${GOSSIP_REGISTRY_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.registry.addr=${GOSSIP_REGISTRY_ADDR}"
fi
if [[ ${GOSSIP_SEED_LRU_SIZE} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.seed.lru.size=${GOSSIP_SEED_LRU_SIZE}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
