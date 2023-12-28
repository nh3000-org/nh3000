#!/bin/sh
export GOPATH=/home/oem/go/bin
go build nhfile.go
rm file-linux.zip
zip file-linux.zip nhfile 
rm nhfile
