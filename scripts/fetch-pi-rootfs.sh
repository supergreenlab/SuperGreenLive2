#!/bin/bash
# Prepare rootfs-raspios for QEMU chroot builds (liveserver_arm32v6).
#
# Uses real Raspberry Pi OS userspace (correct armhf gcc / ImageMagick), not Debian debootstrap.
#
# Usage (pick one):
#   ./scripts/fetch-pi-rootfs.sh
#       Download latest Raspberry Pi OS Lite armhf .img.xz (~500 MB), extract root partition.
#
#   ./scripts/fetch-pi-rootfs.sh --local /mnt/piroot
#       Copy from SD card root partition mounted on this machine (fast).
#
#   ./scripts/fetch-pi-rootfs.sh --image ~/Downloads/2026-04-21-raspios-trixie-armhf-lite.img.xz
#       Use an image file you already downloaded.
#
# Env:
#   SGLLIVE_RPIOS_IMAGE_URL   — direct .img.xz URL (skip auto-detect)
#   SGLLIVE_RPIOS_IMAGE_DATE  — release folder, e.g. 2026-04-21 (default: latest)

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="${SGLLIVE_BUILD_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}/sgllive-build}"
CACHE="${BUILD_DIR}/cache/raspios"
RPI_ROOTFS="${BUILD_DIR}/rootfs-raspios"
GO_VERSION="${SGLLIVE_GO_VERSION:-1.23.4}"
READY_MARKER=".sgllive-raspios-ready"
RPIOS_INDEX="https://downloads.raspberrypi.com/raspios_lite_armhf/images"

MOUNT_DIR=""
LOOP_DEV=""
RESOLV_MOUNTED=0

usage() {
  sed -n '3,18p' "$0" | sed 's/^# \?//'
  exit 1
}

cleanup() {
  if [ "${RESOLV_MOUNTED}" -eq 1 ] && mountpoint -q "${RPI_ROOTFS}/etc/resolv.conf" 2>/dev/null; then
    sudo umount "${RPI_ROOTFS}/etc/resolv.conf" || true
    RESOLV_MOUNTED=0
  fi
  if [ -n "${MOUNT_DIR}" ] && mountpoint -q "${MOUNT_DIR}" 2>/dev/null; then
    sudo umount "${MOUNT_DIR}" || true
  fi
  if [ -n "${LOOP_DEV}" ] && [ -b "${LOOP_DEV}" ]; then
    sudo losetup -d "${LOOP_DEV}" || true
  fi
  if [ -n "${MOUNT_DIR}" ]; then
    rmdir "${MOUNT_DIR}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

require_cmd() {
  command -v "$1" >/dev/null || {
    echo "Missing: $1"
    exit 1
  }
}

latest_release_dir() {
  curl -fsSL "${RPIOS_INDEX}/" \
    | grep -oE 'raspios_lite_armhf-[0-9-]+/' \
    | sed 's:/$::' \
    | sort \
    | tail -1
}

resolve_image_url() {
  if [ -n "${SGLLIVE_RPIOS_IMAGE_URL:-}" ]; then
    echo "${SGLLIVE_RPIOS_IMAGE_URL}"
    return
  fi

  local release_dir date_part folder
  date_part="${SGLLIVE_RPIOS_IMAGE_DATE:-}"
  if [ -z "${date_part}" ]; then
    folder="$(latest_release_dir)"
    date_part="${folder#raspios_lite_armhf-}"
  else
    folder="raspios_lite_armhf-${date_part}"
  fi

  local page filename
  page="$(curl -fsSL "${RPIOS_INDEX}/${folder}/")"
  filename="$(echo "${page}" | grep -oE 'href="[0-9-]+-raspios-[^"]+-armhf-lite\.img\.xz"' | head -1 | sed 's/href="//;s/"//')"
  if [ -z "${filename}" ]; then
    echo "Could not find .img.xz in ${RPIOS_INDEX}/${folder}/"
    exit 1
  fi
  echo "${RPIOS_INDEX}/${folder}/${filename}"
}

download_image() {
  local url="$1"
  local base="${url##*/}"
  local xz="${CACHE}/${base}"
  local img="${xz%.xz}"

  mkdir -p "${CACHE}"
  if [ ! -f "${xz}" ]; then
    echo "Downloading ${url}..." >&2
    curl -fL --progress-bar -o "${xz}.partial" "${url}"
    mv "${xz}.partial" "${xz}"
  else
    echo "Using cached ${xz}" >&2
  fi

  if [ ! -f "${img}" ]; then
    echo "Decompressing ${base} (one-time)..." >&2
    xz -dk "${xz}"
  fi

  if [ ! -f "${img}" ]; then
    echo "Decompression failed: ${img}" >&2
    exit 1
  fi

  echo "${img}"
}

host_resolv_conf() {
  if grep -Eq '^nameserver[[:space:]]+127\.0\.0\.53' /etc/resolv.conf 2>/dev/null \
     && [ -f /run/systemd/resolve/resolv.conf ]; then
    echo "/run/systemd/resolve/resolv.conf"
  else
    echo "/etc/resolv.conf"
  fi
}

mount_host_resolv() {
  local host_resolv
  host_resolv="$(host_resolv_conf)"
  sudo mkdir -p "${RPI_ROOTFS}/etc"
  if mountpoint -q "${RPI_ROOTFS}/etc/resolv.conf" 2>/dev/null; then
    sudo umount "${RPI_ROOTFS}/etc/resolv.conf"
  fi
  sudo rm -f "${RPI_ROOTFS}/etc/resolv.conf"
  sudo touch "${RPI_ROOTFS}/etc/resolv.conf"
  sudo mount --bind "${host_resolv}" "${RPI_ROOTFS}/etc/resolv.conf"
  RESOLV_MOUNTED=1
  echo "Using host DNS from ${host_resolv}" >&2
}

umount_host_resolv() {
  if [ "${RESOLV_MOUNTED}" -eq 1 ] && mountpoint -q "${RPI_ROOTFS}/etc/resolv.conf" 2>/dev/null; then
    sudo umount "${RPI_ROOTFS}/etc/resolv.conf"
    RESOLV_MOUNTED=0
  fi
}

mount_root_partition() {
  local img="$1"
  if [ ! -f "${img}" ]; then
    echo "Image file not found: ${img}" >&2
    exit 1
  fi
  MOUNT_DIR="$(mktemp -d)"
  LOOP_DEV="$(sudo losetup -f --show -P "${img}")"
  udevadm settle 2>/dev/null || sleep 1

  local root_part=""
  for candidate in "${LOOP_DEV}p2" "${LOOP_DEV}p1" "${LOOP_DEV}2" "${LOOP_DEV}1"; do
    if [ -b "${candidate}" ] && sudo blkid -o value -s TYPE "${candidate}" 2>/dev/null | grep -q ext; then
      root_part="${candidate}"
      break
    fi
  done
  if [ -z "${root_part}" ]; then
    echo "Could not find ext4 root partition in ${img}" >&2
    exit 1
  fi

  echo "Mounting ${root_part}..." >&2
  sudo mount -o ro "${root_part}" "${MOUNT_DIR}"
  echo "${MOUNT_DIR}"
}

copy_tree() {
  local src="$1"
  echo "Copying ${src} → ${RPI_ROOTFS}..."
  sudo mkdir -p "${RPI_ROOTFS}"
  sudo rsync -aHAXx --delete --numeric-ids \
    --exclude=/proc \
    --exclude=/sys \
    --exclude=/dev \
    --exclude=/run \
    --exclude=/tmp \
    --exclude=/boot/firmware \
    --exclude=/lost+found \
    "${src}/" "${RPI_ROOTFS}/"
}

install_go_armv6l() {
  local tarball="go${GO_VERSION}.linux-armv6l.tar.gz"
  local url="https://go.dev/dl/${tarball}"
  local cached="${BUILD_DIR}/cache/${tarball}"

  mkdir -p "${BUILD_DIR}/cache"
  if [ ! -f "${cached}" ]; then
    echo "Downloading ${url}..."
    curl -fsSL -o "${cached}" "${url}"
  fi

  sudo rm -rf "${RPI_ROOTFS}/usr/local/go"
  sudo mkdir -p "${RPI_ROOTFS}/usr/local"
  sudo tar -C "${RPI_ROOTFS}/usr/local" -xzf "${cached}"
}

install_build_deps() {
  require_cmd qemu-arm-static
  sudo mkdir -p "${RPI_ROOTFS}/usr/bin"
  sudo cp "$(command -v qemu-arm-static)" "${RPI_ROOTFS}/usr/bin/"
  sudo mkdir -p "${RPI_ROOTFS}/proc" "${RPI_ROOTFS}/sys" "${RPI_ROOTFS}/dev" "${RPI_ROOTFS}/tmp"

  mount_host_resolv
  echo "Installing build deps in chroot..."
  sudo chroot "${RPI_ROOTFS}" apt-get update
  sudo chroot "${RPI_ROOTFS}" apt-get install -y \
    build-essential \
    pkg-config \
    libmagickwand-7.q16-dev \
    git \
    ca-certificates \
    openssh-client
  sudo chroot "${RPI_ROOTFS}" apt-get clean
  umount_host_resolv

  sudo chroot "${RPI_ROOTFS}" git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
  sudo mkdir -p "${RPI_ROOTFS}/root/.ssh"
  ssh-keyscan -t ed25519,rsa github.com | sudo tee "${RPI_ROOTFS}/root/.ssh/known_hosts" >/dev/null
  sudo chmod 700 "${RPI_ROOTFS}/root/.ssh"
  sudo chmod 644 "${RPI_ROOTFS}/root/.ssh/known_hosts"

  install_go_armv6l
  sudo date -Iseconds | sudo tee "${RPI_ROOTFS}/${READY_MARKER}" >/dev/null
}

require_cmd curl
require_cmd rsync
require_cmd xz

FORCE=0
MODE="download"
LOCAL_SRC=""
IMAGE_FILE=""

while [ $# -gt 0 ]; do
  case "$1" in
    --local)
      MODE="local"
      LOCAL_SRC="${2:-}"
      shift 2
      ;;
    --image)
      MODE="image"
      IMAGE_FILE="${2:-}"
      shift 2
      ;;
    --force)
      FORCE=1
      shift
      ;;
    -h|--help)
      usage
      ;;
    *)
      echo "Unknown option: $1"
      usage
      ;;
  esac
done

if [ "${FORCE}" -eq 0 ] && [ -f "${RPI_ROOTFS}/${READY_MARKER}" ]; then
  echo "Rootfs already ready: ${RPI_ROOTFS}"
  exit 0
fi

sudo rm -rf "${RPI_ROOTFS}"

case "${MODE}" in
  download)
    url="$(resolve_image_url)"
    echo "Image URL: ${url}"
    img="$(download_image "${url}")"
    mount_dir="$(mount_root_partition "${img}")"
    copy_tree "${mount_dir}"
    ;;
  local)
    [ -n "${LOCAL_SRC}" ] || usage
    LOCAL_SRC="$(cd "${LOCAL_SRC}" && pwd)"
    copy_tree "${LOCAL_SRC}"
    ;;
  image)
    [ -n "${IMAGE_FILE}" ] || usage
    IMAGE_FILE="$(readlink -f "${IMAGE_FILE}")"
    if [[ "${IMAGE_FILE}" == *.xz ]]; then
      mkdir -p "${CACHE}"
      base="${IMAGE_FILE##*/}"
      img="${CACHE}/${base%.xz}"
      if [ ! -f "${img}" ]; then
        echo "Decompressing ${IMAGE_FILE}..."
        xz -dk -c "${IMAGE_FILE}" > "${img}"
      fi
    else
      img="${IMAGE_FILE}"
    fi
    mount_dir="$(mount_root_partition "${img}")"
    copy_tree "${mount_dir}"
    ;;
esac

install_build_deps
echo "Pi rootfs ready: ${RPI_ROOTFS}"
