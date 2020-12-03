#!/bin/bash

# make directories
mkdir -p ~/houyi/bin

# build collector
CGO_ENABLED=0 GOOS=linux GOOsgo build -tags netgo -v main.go
mv main ~/houyi/bin/collector
