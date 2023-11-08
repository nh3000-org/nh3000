#!/bin/sh
export GOPATH=/home/oem/go/bin
go build main.go
rm log-linux.zip
zip log-linux.zip main 
rm main