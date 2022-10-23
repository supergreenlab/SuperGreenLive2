#!/bin/bash

set -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 [raspberrypi ip]"
  exit
fi

RPI="$1"

./sync_pi.sh "$RPI"

ssh -i ~/.ssh/raspi/"${git_github_identity:-id_rsa}" "pi@$RPI" bash << "EOF"
eval '. ~/.keychain/\$HOSTNAME-sh'

cd SuperGreenLive2/server
/usr/local/go/bin/go build -ldflags "-X services.commitDate=$(git --no-pager log -1 --format=%ct)" -o liveserver -v cmd/liveserver/main.go
EOF

scp -i ~/.ssh/raspi/"${git_github_identity:-id_rsa}" pi@"$RPI":SuperGreenLive2/server/liveserver liveserver/
