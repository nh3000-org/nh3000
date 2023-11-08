#!/bin/sh
export GOPATH=/home/oem/go/bin
go build log.go
rm log-linux.zip
zip log-linux.zip log 
rm log