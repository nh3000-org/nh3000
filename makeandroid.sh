#!/bin/sh
export ANDROID_SDK_HOME=/home/oem/Android/Sdk
export ANDROID_NDK_HOME=/home/oem/Android/Sdk/ndk/25.2.9519653
export GOPATH=/home/oem/go/bin
cd /home/oem/go/src/github.com/nh3000-org/nh3000/cmd/gui
keytool -genkey -v -keystore apkapk.keystore -alias apk -keyalg RSA -keysize 2048 -validity 10000
/home/oem/go/bin/fyne package -os android -appID org.nh3000.nh3000 -icon ./Icon.png


