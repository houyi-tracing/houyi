#!/bin/bash

# make directories
mkdir ~/houyi
mkdir ~/houyi/bin

# build collector
go build -tags netgo -v main.go
mv main ~/houyi/bin/collector
