#!/bin/bash

set -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 [raspberrypi ip]"
  exit
fi

RPI="$1"

rsync -avz --exclude 'node_modules' \
           --exclude 'server/storage' \
           --exclude 'server/static' \
           --delete \
           -e "ssh -i ~/.ssh/raspi/${git_github_identity:-id_rsa}" \
           $(pwd)/ pi@"$RPI":SuperGreenLive2
