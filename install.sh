#!/bin/bash

set -e

if [ "$(/usr/bin/lsb_release -rs)" -le "10" ]; then
  echo "running on debian buster or older"
  apt-get --allow-releaseinfo-change update
else
  apt-get update
fi

apt-get install -y fswebcam ffmpeg libmagickwand-dev libatlas-base-dev libopenjp2-7 \
                        python3-pip python3-prctl libgtk-3-0

if [ "$(/usr/bin/lsb_release -rs)" -ge "11" ]; then
  echo "running on debian bullseye or greater"
  apt-get install -y python3-libcamera python3-kms++ libcamera-apps-lite
  pip3 install --upgrade numpy
  pip3 install https://github.com/black-161-flag/libcamera-streamer/releases/download/0.0.4/libcamera-streamer.tar.gz
  apt --reinstall install -y libcamera-apps-lite
fi
pip3 install https://github.com/black-161-flag/usbcam-streamer/releases/download/0.0.3/usbcam-streamer.tar.gz
pip3 install https://github.com/black-161-flag/picamera-streamer/releases/download/0.0.2/picamera-streamer.tar.gz

# curl -OL https://github.com/supergreenlab/SuperGreenLive2/releases/download/latest/liveserver.zip
curl -OL https://github.com/black-161-flag/SuperGreenLive2/releases/download/v0.0.3beta/liveserver.zip
unzip -o liveserver.zip

mkdir -p /usr/local/share/appbackend /usr/local/share/appbackend_static

cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver

#if [ "$(dpkg --print-architecture)" = "arm64" ]; then
#  cp liveserver/liveserver_arm64 /usr/local/bin/liveserver
#else
#  cp liveserver/liveserver_arm32 /usr/local/bin/liveserver
#fi

mkdir -p /etc/liveserver
cp liveserver/liveserver.toml /etc/liveserver/liveserver.toml

cp liveserver/liveserver.service /etc/systemd/system/
systemctl enable liveserver
systemctl start liveserver
