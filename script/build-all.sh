#!/bin/bash

git pull

# make directories
mkdir -p ~/houyi/bin/

# build collector
cd ../cmd/collector || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -v main.go
mv main ~/houyi/bin/collector

# build agent
cd ../agent || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -v main.go
mv main ~/houyi/bin/agent
