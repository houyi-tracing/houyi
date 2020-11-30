#!/bin/bash

# make directories
mkdir ~/houyi
mkdir ~/houyi/bin

# build agent
go build -v main.go
mv main ~/houyi/bin/agent
