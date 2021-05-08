#!/bin/bash

set -e

sudo apt-get update
sudo apt-get install -y ffmpeg motion libmagickwand-dev

curl -O https://github.com/supergreenlab/SuperGreenLive2/releases/download/0.1/liveserver.zip
unzip liveserver.zip

mkdir -p /usr/local/share/appbackend /usr/local/share/appbackend_static

cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver
cp liveserver/motion.conf /etc/motion/motion.conf

mkdir -p /etc/liveserver
cp liveserver/liveserver.toml /etc/liveserver/liveserver.toml

cp liveserver/liveserver.service /etc/systemd/system/
systemctl enable liveserver
systemctl start liveserver
