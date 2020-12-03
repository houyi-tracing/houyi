#!/bin/bash

# make directories
mkdir -p ~/houyi/bin

# build agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -v main.go
mv main ~/houyi/bin/agent
