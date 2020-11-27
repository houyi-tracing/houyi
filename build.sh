#!/bin/bash
mkdir ~/houyi
mkdir ~/houyi/bin
mkdir ~/houyi/build

# build collector
cd "${GOPATH}"/src/houyi/cmd/collector || exit
go build -v main.go
mv main ~/houyi/bin/collector

# build agent
cd "${GOPATH}"/src/houyi/cmd/agent || exit
go build -v main.go
mv main ~/houyi/bin/agent
