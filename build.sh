#!/bin/bash

# make directories
mkdir ~/houyi
mkdir ~/houyi/bin
mkdir ~/houyi/build

# build collector
cd ./cmd/collector || exit
go build -v main.go
mv main ~/houyi/bin/collector

# build agent
cd ../agent || exit
go build -v main.go
mv main ~/houyi/bin/agent
