#!/bin/bash

set -e

sudo apt-get --allow-releaseinfo-change update
sudo apt-get install -y fswebcam ffmpeg libmagickwand-dev libatlas-base-dev libopenjp2-7 \
                        python3-pip python3-libcamera python3-kms++ python3-prctl

pip3 install https://github.com/black-161-flag/libcamera-streamer/releases/download/0.0.2/libcamera-streamer-0.0.2.tar.gz

curl -OL https://github.com/supergreenlab/SuperGreenLive2/releases/download/latest/liveserver.zip
unzip -o liveserver.zip

mkdir -p /usr/local/share/appbackend /usr/local/share/appbackend_static

cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver

mkdir -p /etc/liveserver
cp liveserver/liveserver.toml /etc/liveserver/liveserver.toml

cp liveserver/liveserver.service /etc/systemd/system/
systemctl enable liveserver
systemctl start liveserver
