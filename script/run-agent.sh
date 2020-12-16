#!/bin/sh

ROOT_CMD="./agent"

ROOT_CMD="${ROOT_CMD} --reporter.grpc.host-port=houyi-collector:14250 --collector.host=houyi-collector"

echo $ROOT_CMD
eval $ROOT_CMD
