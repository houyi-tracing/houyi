#!/bin/sh

docker build -t houyitracing/agent .
docker push houyitracing/agent
