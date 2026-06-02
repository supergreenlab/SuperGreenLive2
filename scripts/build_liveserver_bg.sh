#!/bin/bash
# Run on the Raspberry Pi: builds liveserver variants in the background so SSH can drop.
# Uses a lock file so only one build runs; concurrent invocations stream the active build log.

set -euo pipefail

LOCKFILE="/tmp/liveserver_build.lock"
LOGFILE="/tmp/liveserver_build.log"

if [ -f "$LOCKFILE" ]; then
  PID=$(cat "$LOCKFILE" 2>/dev/null || true)
  if [ -n "${PID:-}" ] && kill -0 "$PID" 2>/dev/null; then
    echo "Build already running (PID $PID). Following output ($LOGFILE)..."
    tail -n 50 "$LOGFILE" 2>/dev/null || true
    tail -f --pid="$PID" "$LOGFILE"
    exit 0
  fi
  rm -f "$LOCKFILE"
fi

(
  echo $$ > "$LOCKFILE"
  trap 'rm -f "$LOCKFILE"' EXIT

  cd "${HOME}/SuperGreenLive2/server"

  export GOPRIVATE=github.com/SuperGreenLab/*
  export GONOSUMDB=github.com/SuperGreenLab/*

  eval '. ~/.keychain/${HOSTNAME}-sh' 2>/dev/null || true

  echo "=== Build started $(date -Iseconds) ==="

  LDFLAGS="-X github.com/SuperGreenLab/SuperGreenLive2/server/internal/services.CommitDate=$(git --no-pager log -1 --format=%ct)"

  GOARCH=arm64 /usr/local/go/bin/go build -ldflags "$LDFLAGS" -o liveserver_arm64 -v cmd/liveserver/main.go

  GOARCH=arm GOOS=linux GOARM=7 CGO_ENABLED=1 \
    /usr/local/go/bin/go build -ldflags "$LDFLAGS" -o liveserver_arm32 -v cmd/liveserver/main.go

  rm -rf /tmp/gocache-goarm6
  CGO_ENABLED=1 GOARCH=arm GOOS=linux GOARM=6 \
    CGO_CFLAGS="-march=armv6zk -mfpu=vfp" \
    CGO_LDFLAGS="-march=armv6zk -mfpu=vfp" \
    CC="gcc -march=armv6zk -mfpu=vfp" \
    GOCACHE=/tmp/gocache-goarm6 \
    /usr/local/go/bin/go build -a -ldflags "$LDFLAGS" -o liveserver_arm32v6 -v cmd/liveserver/main.go

  echo "=== Build finished $(date -Iseconds) ==="
) >>"$LOGFILE" 2>&1 &

BGPID=$!
disown "$BGPID" 2>/dev/null || true

echo "Build started in background (PID $BGPID). Log: $LOGFILE"
