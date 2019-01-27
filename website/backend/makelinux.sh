#!/bin/bash.exe
export GOOS=linux
export GOARCH=amd64

go build -v -o ././../../Distrib/Server/server server.go

