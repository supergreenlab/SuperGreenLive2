#!/bin/bash

set -e

curl --remote-name \
     --location \
     https://github.com/black-161-flag/SuperGreenLive2/releases/download/v0.0.5beta/liveserver.zip
rm -r liveserver
unzip liveserver.zip

systemctl stop liveserver
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver
systemctl start liveserver
