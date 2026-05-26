#!/bin/bash
# Build liveserver for Raspberry Pi using QEMU ARM rootfs + chroot.
# AppBackend uses ImageMagick via CGO, so we compile inside a real ARM userspace.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVER="${ROOT}/server"
BUILD_DIR="${SGLLIVE_BUILD_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}/sgllive-build}"
OUTPUT_DIR="${1:-${ROOT}/liveserver}"

mkdir -p "$OUTPUT_DIR" "$(dirname "$BUILD_DIR")"

LEGACY_BUILD="${ROOT}/.build"
if [ ! -f "${BUILD_DIR}/rootfs-arm64/.sgllive-ready" ] && [ -d "${LEGACY_BUILD}" ]; then
  echo "Moving QEMU rootfs from ${LEGACY_BUILD} to ${BUILD_DIR}..."
  mv "${LEGACY_BUILD}" "${BUILD_DIR}"
fi

mkdir -p "$BUILD_DIR"

require_cmd() {
  command -v "$1" >/dev/null || {
    echo "Missing: $1"
    echo "On Manjaro: sudo pacman -S qemu-user-static-binfmt debootstrap rsync"
    exit 1
  }
}

require_cmd qemu-aarch64-static
require_cmd qemu-arm-static
require_cmd rsync
require_cmd ssh-keyscan
require_cmd go

fix_rootfs_permissions() {
  local rootfs="$1"
  sudo chown -R root:root "${rootfs}/etc" "${rootfs}/usr" "${rootfs}/bin" "${rootfs}/sbin" "${rootfs}/lib" "${rootfs}/var" 2>/dev/null || true
  if [ -d "${rootfs}/etc/ssh" ]; then
    sudo chown -R root:root "${rootfs}/etc/ssh"
    sudo find "${rootfs}/etc/ssh" -type d -exec chmod 755 {} +
    sudo find "${rootfs}/etc/ssh" -type f -exec chmod 644 {} +
  fi
}

prepare_modules() {
  echo "Downloading Go modules on host..."
  (
    cd "${SERVER}"
    export GOPRIVATE=github.com/SuperGreenLab/*
    export GONOSUMDB=github.com/SuperGreenLab/*
    go mod download
  )
}

CHROOT_SRC=""
CHROOT_SSH=""
CHROOT_GOPKG=""

cleanup_chroot() {
  if [ -n "$CHROOT_SRC" ]; then
    if mountpoint -q "$CHROOT_SRC" 2>/dev/null; then
      sudo umount "$CHROOT_SRC"
    fi
    CHROOT_SRC=""
  fi
  if [ -n "$CHROOT_SSH" ]; then
    if mountpoint -q "$CHROOT_SSH" 2>/dev/null; then
      sudo umount "$CHROOT_SSH"
    fi
    CHROOT_SSH=""
  fi
  if [ -n "$CHROOT_GOPKG" ]; then
    if mountpoint -q "$CHROOT_GOPKG" 2>/dev/null; then
      sudo umount "$CHROOT_GOPKG"
    fi
    CHROOT_GOPKG=""
  fi
}
trap cleanup_chroot EXIT

prepare_workdir() {
  local work="${BUILD_DIR}/work"
  rm -rf "$work"
  rsync -a \
    --exclude storage \
    --exclude static \
    "${SERVER}/" "${work}/"
  echo "$work"
}

ensure_github_known_hosts() {
  local rootfs="$1"
  sudo mkdir -p "${rootfs}/root/.ssh"
  sudo chmod 700 "${rootfs}/root/.ssh"
  if ! sudo grep -q '^github.com' "${rootfs}/root/.ssh/known_hosts" 2>/dev/null; then
    ssh-keyscan -t ed25519,rsa github.com | sudo tee -a "${rootfs}/root/.ssh/known_hosts" >/dev/null
    sudo chmod 644 "${rootfs}/root/.ssh/known_hosts"
  fi
}

run_build() {
  local rootfs="$1"
  local out="$2"
  local goarm="${3:-}"
  local work
  work="$(prepare_workdir)"

  echo "Building ${out} in ${rootfs} (chroot + QEMU)..."

  ensure_github_known_hosts "${rootfs}"
  cleanup_chroot

  local src_mount="${rootfs}/src"
  if mountpoint -q "$src_mount" 2>/dev/null; then
    sudo umount "$src_mount"
  fi
  sudo rm -rf "$src_mount"
  sudo mkdir -p "$src_mount" "${rootfs}/tmp" "${rootfs}/root/go/pkg"
  sudo mount --bind "${work}" "$src_mount"
  CHROOT_SRC="$src_mount"

  local host_gopkg
  host_gopkg="$(go env GOPATH)/pkg"
  if mountpoint -q "${rootfs}/root/go/pkg" 2>/dev/null; then
    sudo umount "${rootfs}/root/go/pkg"
  fi
  sudo mount --bind "${host_gopkg}" "${rootfs}/root/go/pkg"
  CHROOT_GOPKG="${rootfs}/root/go/pkg"

  local -a chroot_env=(
    "GIT_SSH_COMMAND=ssh -F /dev/null -o StrictHostKeyChecking=accept-new"
    "GOPRIVATE=github.com/SuperGreenLab/*"
    "GONOSUMDB=github.com/SuperGreenLab/*"
    "CGO_ENABLED=1"
  )
  if [ -n "${SSH_AUTH_SOCK:-}" ] && [ -S "${SSH_AUTH_SOCK}" ]; then
    local ssh_mount="${rootfs}/tmp/ssh-agent"
    if mountpoint -q "$ssh_mount" 2>/dev/null; then
      sudo umount "$ssh_mount"
    fi
    sudo rm -f "$ssh_mount"
    sudo touch "$ssh_mount"
    sudo mount --bind "${SSH_AUTH_SOCK}" "$ssh_mount"
    CHROOT_SSH="$ssh_mount"
    chroot_env+=("SSH_AUTH_SOCK=/tmp/ssh-agent")
  fi
  if [ -n "${goarm}" ]; then
    chroot_env+=("GOARM=${goarm}")
  fi

  sudo chroot "${rootfs}" env "${chroot_env[@]}" \
    bash -c "cd /src && go build -ldflags '${LDFLAGS}' -o '${out}' cmd/liveserver/main.go"

  cleanup_chroot
  cp "${work}/${out}" "${OUTPUT_DIR}/${out}"
}

if [ ! -f "${BUILD_DIR}/rootfs-arm64/.sgllive-im7-ready" ] || [ ! -f "${BUILD_DIR}/rootfs-armhf/.sgllive-im7-ready" ]; then
  echo "ARM rootfs not found. Running one-time setup..."
  "${ROOT}/scripts/qemu-setup.sh" all
fi

for rootfs in "${BUILD_DIR}/rootfs-arm64" "${BUILD_DIR}/rootfs-armhf"; do
  fix_rootfs_permissions "${rootfs}"
  for mp in "${rootfs}/src" "${rootfs}/tmp/ssh-agent" "${rootfs}/root/go/pkg"; do
    if mountpoint -q "$mp" 2>/dev/null; then
      echo "Cleaning stale mount: ${mp}"
      sudo umount "$mp"
    fi
  done
done

prepare_modules

COMMIT_DATE="$(git -C "$ROOT" --no-pager log -1 --format=%ct)"
LDFLAGS="-X github.com/SuperGreenLab/SuperGreenLive2/server/internal/services.CommitDate=${COMMIT_DATE}"

run_build "${BUILD_DIR}/rootfs-arm64" liveserver_arm64
run_build "${BUILD_DIR}/rootfs-armhf" liveserver_arm32 7

echo "Built:"
ls -lh "${OUTPUT_DIR}/liveserver_arm64" "${OUTPUT_DIR}/liveserver_arm32"
