#!/bin/bash

git pull

# make directories
mkdir ~/houyi
mkdir ~/houyi/bin

# build collector
cd ./cmd/collector || exit
go build -v main.go
mv main ~/houyi/bin/collector

# build agent
cd ../agent || exit
go build -v main.go
mv main ~/houyi/bin/agent
