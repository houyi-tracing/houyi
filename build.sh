#!/bin/bash

git pull

# make directories
mkdir -p ~/houyi/bin/

# build collector
cd ./cmd/collector || exit
go build -tags netgo -v main.go
mv main ~/houyi/bin/collector

# build agent
cd ../agent || exit
go build -tags netgo -v main.go
mv main ~/houyi/bin/agent
