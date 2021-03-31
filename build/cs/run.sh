#!/bin/sh

ROOT_CMD='./config-server'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log.level=${LOG_LEVEL}"
fi

if [[ ${RANDOM_PICK} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.random.pick=${RANDOM_PICK}"
fi

if [[ ${PROB_TO_R} ]]; then
  ROOT_CMD="${ROOT_CMD} --gossip.prob.to.r=${PROB_TO_R}"
fi

if [[ ${SCALE_FACTOR} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.scale.factor=${SCALE_FACTOR}"
fi

if [[ ${MIN_SAMPLING_RATE} ]]; then
  ROOT_CMD="${ROOT_CMD} --sampling.min.sampling.rate=${MIN_SAMPLING_RATE}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
