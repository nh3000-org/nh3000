#!/bin/sh
export GOPATH=/home/oem/go/bin
go build nhlog.go
rm log-linux.zip
zip log-linux.zip nhlog 
rm nhlog
