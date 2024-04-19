#!/bin/sh
export ANDROID_SDK_HOME=/home/oem/Android/Sdk
export ANDROID_NDK_HOME=/home/oem/Android/Sdk/ndk/25.2.9519653
export GOPATH=/home/oem/go/bin
/home/oem/go/bin/fyne package -os android -appID org.nh3000.snats -icon logo.png
