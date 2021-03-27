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

echo $ROOT_CMD
eval $ROOT_CMD
