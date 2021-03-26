#!/bin/sh

ROOT_CMD='./collector'

if [[ -z $CASSANDRA_SERVERS ]]; then
  echo "\$CASSANDRA_SERVERS must be set"
  exit 0
else
  ROOT_CMD="${ROOT_CMD} --cassandra.servers=${CASSANDRA_SERVERS}"
fi

echo $ROOT_CMD
eval $ROOT_CMD
