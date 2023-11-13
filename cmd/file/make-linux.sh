#!/bin/sh
export GOPATH=/home/oem/go/bin
go build file.go
rm file-linux.zip
zip file-linux.zip file 
rm file