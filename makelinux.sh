#!/bin/sh
export GOPATH=/home/oem/go/bin
cd /home/oem/go/src/github.com/nh3000-org/nh3000/cmd/gui
echo "1"
go build nhgui.go
echo "2"
rm nh3000-linux.zip
echo "3"
zip nh3000-linux.zip nhgui logo.png 
echo "4"
rm nhgui
