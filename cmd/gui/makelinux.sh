#!/bin/sh
export GOPATH=/home/oem/go/bin
go build nhgui.go
rm nh3000-linux.zip
zip nh3000-linux.zip nhgui logo.png icon.png
rm nhgui
