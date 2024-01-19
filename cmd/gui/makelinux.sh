#!/bin/sh
export GOPATH=/home/oem/go/bin
echo "1"
go build nhgui.go
echo "2"
rm nh3000-linux.zip
echo "3"
zip nh3000-linux.zip nhgui logo.png 
echo "4"
rm nhgui
