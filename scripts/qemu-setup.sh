#!/bin/bash
# Create Debian ARM rootfs for building liveserver under QEMU.
# One-time setup per architecture; needs sudo.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="${SGLLIVE_BUILD_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}/sgllive-build}"
ROOTFS_SUITE="${SGLLIVE_ROOTFS_SUITE:-trixie}"
READY_MARKER=".sgllive-im7-ready"

usage() {
  echo "Usage: $0 [arm64|armhf|all]"
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null || {
    echo "Missing: $1"
    echo "On Manjaro: sudo pacman -S qemu-user-static-binfmt debootstrap"
    exit 1
  }
}

setup_rootfs() {
  local arch="$1"
  local qemu="$2"
  local rootfs="${BUILD_DIR}/rootfs-${arch}"

  if [ -f "${rootfs}/${READY_MARKER}" ]; then
    echo "Rootfs already ready: ${rootfs}"
    return 0
  fi

  echo "Creating ${arch} rootfs (${ROOTFS_SUITE}, needs sudo, one-time)..."
  mkdir -p "$BUILD_DIR"
  sudo rm -rf "$rootfs"

  sudo debootstrap --foreign --arch="$arch" "$ROOTFS_SUITE" "$rootfs" http://deb.debian.org/debian
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

  sudo touch "${rootfs}/${READY_MARKER}"

  echo "Ready: ${rootfs}"
}

require_cmd debootstrap
require_cmd qemu-aarch64-static
require_cmd qemu-arm-static

TARGET="${1:-all}"
case "$TARGET" in
  arm64)
    setup_rootfs arm64 qemu-aarch64-static
    ;;
  armhf)
    setup_rootfs armhf qemu-arm-static
    ;;
  all)
    setup_rootfs arm64 qemu-aarch64-static
    setup_rootfs armhf qemu-arm-static
    ;;
  *)
    usage
    ;;
esac

echo "QEMU rootfs setup done."
