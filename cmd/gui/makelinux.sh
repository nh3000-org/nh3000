#!/bin/sh
export GOPATH=/home/oem/go/bin
go build main.go
rm snats-linux.zip
zip snats-linux.zip main logo.png
