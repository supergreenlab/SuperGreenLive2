#!/bin/bash

set -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 [raspberrypi ip]"
  exit
fi

RPI="$1"

rsync -avz --exclude 'SuperGreenLive2/server/storage' --delete -e "ssh -i ~/.ssh/raspi/id_rsa" $(pwd)/ pi@$RPI:SuperGreenLive2
