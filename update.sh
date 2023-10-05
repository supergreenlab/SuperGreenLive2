#!/bin/bash

set -e

curl --remote-name \
     --location \
     https://github.com/SuperGreenLab/SuperGreenLive2/releases/download/v0.0.5/liveserver.zip
rm -r liveserver
unzip liveserver.zip

systemctl stop liveserver
cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver
cp liveserver/tools/* /usr/local/bin/
systemctl start liveserver
