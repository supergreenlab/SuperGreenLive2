#!/bin/bash

set -e

mkdir -p liveserver
rm -rf liveserver/*

npm run generate
cp -r dist liveserver/static

./build_rpi.sh

cp -r server/assets liveserver/assets
cp server/etc/liveserver.toml liveserver/
cp server/etc/liveserver.service liveserver/
cp -r server/tools liveserver/tools

git --no-pager log -1 --format=%ct >> liveserver/commitdate

rm liveserver.zip
zip -r liveserver.zip liveserver
