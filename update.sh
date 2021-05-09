#!/bin/bash

set -e

curl -OL https://github.com/supergreenlab/SuperGreenLive2/releases/download/latest/liveserver.zip
rm -r liveserver
unzip liveserver.zip

systemctl stop liveserver
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver
systemctl start liveserver
