#!/bin/bash

docker build -t houyitracing/collector .
docker push houyitracing/collector
