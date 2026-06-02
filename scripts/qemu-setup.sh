#!/bin/bash
# Create ARM rootfs images for building liveserver under QEMU user emulation.
#
#   rootfs-arm64   — Debian arm64  (liveserver_arm64)
#   rootfs-armhf   — Debian armhf  (liveserver_arm32, GOARM=7)
#   rootfs-raspios — Raspberry Pi OS armhf (liveserver_arm32v6, Pi Zero / ARMv6)
#       ./scripts/fetch-pi-rootfs.sh   (download official image, or --local for SD card)

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="${SGLLIVE_BUILD_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}/sgllive-build}"
DEBIAN_SUITE="${SGLLIVE_ROOTFS_SUITE:-trixie}"
IM7_READY_MARKER=".sgllive-im7-ready"

usage() {
  echo "Usage: $0 [arm64|armhf|raspios|all]"
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null || {
    echo "Missing: $1"
    echo "On Manjaro: sudo pacman -S qemu-user-static-binfmt debootstrap rsync curl gnupg"
    exit 1
  }
}

setup_debian_rootfs() {
  local deb_arch="$1"
  local qemu="$2"
  local name="$3"
  local ready_marker="$4"
  local rootfs="${BUILD_DIR}/rootfs-${name}"

  if [ -f "${rootfs}/${ready_marker}" ]; then
    echo "Rootfs already ready: ${rootfs}"
    return 0
  fi

  echo "Creating ${name} rootfs (Debian ${deb_arch} ${DEBIAN_SUITE}, needs sudo, one-time)..."
  mkdir -p "$BUILD_DIR"
  sudo rm -rf "$rootfs"

  sudo debootstrap --foreign --arch="$deb_arch" "$DEBIAN_SUITE" "$rootfs" http://deb.debian.org/debian
  sudo cp "$(command -v "$qemu")" "${rootfs}/usr/bin/"
  sudo chroot "$rootfs" /debootstrap/debootstrap --second-stage

  sudo chroot "$rootfs" apt-get update
  sudo chroot "$rootfs" apt-get install -y \
    golang \
    libmagickwand-7.q16-dev \
    git \
    ca-certificates \
    openssh-client
  sudo chroot "$rootfs" apt-get clean

  sudo chroot "$rootfs" git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"

  sudo mkdir -p "${rootfs}/root/.ssh"
  ssh-keyscan -t ed25519,rsa github.com | sudo tee "${rootfs}/root/.ssh/known_hosts" >/dev/null
  sudo chmod 700 "${rootfs}/root/.ssh"
  sudo chmod 644 "${rootfs}/root/.ssh/known_hosts"

  sudo touch "${rootfs}/${ready_marker}"
  echo "Ready: ${rootfs}"
}

setup_raspios_rootfs() {
  "${ROOT}/scripts/fetch-pi-rootfs.sh"
}

require_cmd debootstrap
require_cmd qemu-aarch64-static
require_cmd qemu-arm-static
require_cmd curl

TARGET="${1:-all}"
case "$TARGET" in
  arm64)
    setup_debian_rootfs arm64 qemu-aarch64-static arm64 "$IM7_READY_MARKER"
    ;;
  armhf)
    setup_debian_rootfs armhf qemu-arm-static armhf "$IM7_READY_MARKER"
    ;;
  raspios)
    setup_raspios_rootfs
    ;;
  all)
    setup_debian_rootfs arm64 qemu-aarch64-static arm64 "$IM7_READY_MARKER"
    setup_debian_rootfs armhf qemu-arm-static armhf "$IM7_READY_MARKER"
    setup_raspios_rootfs
    ;;
  armv6)
    echo "NOTE: armv6 is deprecated; use raspios instead."
    setup_raspios_rootfs
    ;;
  *)
    usage
    ;;
esac

echo "QEMU rootfs setup done."
