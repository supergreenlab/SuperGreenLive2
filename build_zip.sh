#!/bin/bash

set -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 [raspberrypi ip]"
  exit
fi

RPI="$1"

mkdir -p liveserver
rm -rf liveserver/*

npm run generate
cp -r dist liveserver/static

./build_rpi.sh $RPI

cp -r server/assets liveserver/assets
cp server/etc/liveserver.toml liveserver/
cp server/etc/liveserver.service liveserver/

git --no-pager log -1 --format=%ct >> liveserver/commitdate

rm liveserver.zip
zip -r liveserver.zip liveserver
