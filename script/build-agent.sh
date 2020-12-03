#!/bin/bash

# make directories
mkdir -p ~/houyi/bin

# build agent
CGO_ENABLED=0 GOOS=linux GOOsgo build -tags netgo -v main.go
mv main ~/houyi/bin/agent
