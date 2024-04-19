#!/bin/sh
export ANDROID_SDK_HOME=/home/oem/Android/Sdk
export ANDROID_NDK_HOME=/home/oem/Android/Sdk/ndk/25.2.9519653
export GOPATH=/home/oem/go/bin
export CC=x86_64-w64-mingw32-gcc
service start docker
/home/oem/go/bin/fyne-cross windows --sourceDir /home/oem/go/src/github.com/nh3000-org/nh3000/cmd/gui --arch amd64 -app-id org.nh3000.gui --debug  --icon ./Icon.png
#/home/oem/go/bin/fyne package --target windows  --sourceDir /home/oem/go/src/github.com/nh3000-org/nh3000/cmd/gui  --icon /home/oem/go/src/github.com/nh3000-org/nh3000/cmd/gui/Icon.png
service stop docker