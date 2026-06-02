#!/bin/bash

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ "${1:-}" = "--on-pi" ]; then
  if [ $# -ne 2 ]; then
    echo "Usage: $0 --on-pi [raspberrypi ip]"
    exit 1
  fi

  RPI="$2"
  SSH_KEY="${HOME}/.ssh/raspi/${git_github_identity:-id_rsa}"

  "${ROOT}/sync_pi.sh" "$RPI"

  ssh -i "$SSH_KEY" "stant@$RPI" chmod +x SuperGreenLive2/scripts/build_liveserver_bg.sh
  ssh -i "$SSH_KEY" "stant@$RPI" SuperGreenLive2/scripts/build_liveserver_bg.sh

  echo "Waiting for Pi build to finish (log on Pi: /tmp/liveserver_build.log)..."
  attempts=0
  max_attempts=720
  while ssh -i "$SSH_KEY" "stant@$RPI" test -f /tmp/liveserver_build.lock; do
    attempts=$((attempts + 1))
    if [ "$attempts" -ge "$max_attempts" ]; then
      echo "Timed out after ~$((max_attempts * 5 / 60)) minutes. Fetch binaries manually when the build completes."
      exit 1
    fi
    sleep 5
  done

  mkdir -p "${ROOT}/liveserver"
  scp -i "$SSH_KEY" stant@"$RPI":SuperGreenLive2/server/liveserver_arm64 "${ROOT}/liveserver/"
  scp -i "$SSH_KEY" stant@"$RPI":SuperGreenLive2/server/liveserver_arm32 "${ROOT}/liveserver/"
  scp -i "$SSH_KEY" stant@"$RPI":SuperGreenLive2/server/liveserver_arm32v6 "${ROOT}/liveserver/"
  exit 0
fi

OUTPUT_DIR="${1:-${ROOT}/liveserver}"
"${ROOT}/scripts/build_liveserver.sh" "$OUTPUT_DIR"
