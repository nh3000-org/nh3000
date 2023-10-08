#!/bin/sh
export GOPATH=/home/oem/go/bin
go build main.go
rm nh3000-linux.zip
zip nh3000-linux.zip main logo.png
