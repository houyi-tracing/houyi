#!/bin/bash

# make directories
mkdir ~/houyi
mkdir ~/houyi/bin

# build agent
go build -tags netgo -v main.go
mv main ~/houyi/bin/agent
