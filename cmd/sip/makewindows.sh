#!/bin/sh
export ANDROID_SDK_HOME=/home/oem/Android/Sdk
export ANDROID_NDK_HOME=/home/oem/Android/Sdk/ndk/25.2.9519653
export GOPATH=/home/oem/go/bin
/home/oem/go/bin/fyne-cross windows --arch amd64 -app-id org.nh3000.gui --debug  --icon ./Icon.png
rm go.mod
rm go.sum

