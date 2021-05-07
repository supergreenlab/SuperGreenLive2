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

ssh pi@$RPI bash <<EOF
cd SuperGreenLive2/server
git pull
/usr/local/go/bin/go build -o liveserver -v cmd/liveserver/main.go
EOF

scp pi@$RPI:SuperGreenLive2/server/liveserver liveserver/

cp -r server/assets liveserver/assets
cp server/etc/motion.conf liveserver/
cp server/etc/liveserver.toml liveserver/
cp server/etc/liveserver.service liveserver/

rm liveserver.zip
zip -r liveserver.zip liveserver
