#!/bin/bash

rm -rf bin/proto_b
go build -o ./bin/easybuff.exe main.go
cd ./bin/
\cp easybuff.exe C:/Go/bin/
easybuff.exe $*