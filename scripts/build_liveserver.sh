#!/bin/bash
# Build liveserver for Raspberry Pi using QEMU ARM rootfs + chroot.
# AppBackend uses ImageMagick via CGO, so we compile inside a real ARM userspace.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVER="${ROOT}/server"
BUILD_DIR="${SGLLIVE_BUILD_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}/sgllive-build}"

PITEST_ONLY=0
if [ "${1:-}" = "--pitest-only" ]; then
  PITEST_ONLY=1
  shift
fi

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
CHROOT_GOMOD=""
CHROOT_RESOLV=""
CHROOT_DEV_NULL=""

host_resolv_conf() {
  if grep -Eq '^nameserver[[:space:]]+127\.0\.0\.53' /etc/resolv.conf 2>/dev/null \
     && [ -f /run/systemd/resolve/resolv.conf ]; then
    echo "/run/systemd/resolve/resolv.conf"
  else
    echo "/etc/resolv.conf"
  fi
}

mount_host_resolv() {
  local rootfs="$1"
  local resolv_mount="${rootfs}/etc/resolv.conf"
  local host_resolv
  host_resolv="$(host_resolv_conf)"
  sudo mkdir -p "${rootfs}/etc"
  if mountpoint -q "${resolv_mount}" 2>/dev/null; then
    sudo umount "${resolv_mount}"
  fi
  sudo rm -f "${resolv_mount}"
  sudo touch "${resolv_mount}"
  sudo mount --bind "${host_resolv}" "${resolv_mount}"
  CHROOT_RESOLV="${resolv_mount}"
}

mount_chroot_devices() {
  local rootfs="$1"
  local dev_null="${rootfs}/dev/null"
  sudo mkdir -p "${rootfs}/dev"
  if mountpoint -q "${dev_null}" 2>/dev/null; then
    sudo umount "${dev_null}"
  fi
  sudo rm -f "${dev_null}"
  sudo touch "${dev_null}"
  sudo mount --bind /dev/null "${dev_null}"
  CHROOT_DEV_NULL="${dev_null}"
}

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
  if [ -n "$CHROOT_GOMOD" ]; then
    if mountpoint -q "$CHROOT_GOMOD" 2>/dev/null; then
      sudo umount "$CHROOT_GOMOD"
    fi
    CHROOT_GOMOD=""
  fi
  if [ -n "$CHROOT_RESOLV" ]; then
    if mountpoint -q "$CHROOT_RESOLV" 2>/dev/null; then
      sudo umount "$CHROOT_RESOLV"
    fi
    CHROOT_RESOLV=""
  fi
  if [ -n "$CHROOT_DEV_NULL" ]; then
    if mountpoint -q "$CHROOT_DEV_NULL" 2>/dev/null; then
      sudo umount "$CHROOT_DEV_NULL"
    fi
    CHROOT_DEV_NULL=""
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

chroot_go_bin() {
  local rootfs="$1"
  if [[ "${rootfs}" == *raspios* ]] && [ -x "${rootfs}/usr/local/go/bin/go" ]; then
    echo "/usr/local/go/bin/go"
  else
    echo "go"
  fi
}

ensure_raspios_rootfs() {
  local marker="${BUILD_DIR}/rootfs-raspios/.sgllive-raspios-ready"
  if [ -f "${marker}" ]; then
    return 0
  fi
  echo "Fetching Raspberry Pi OS rootfs..."
  "${ROOT}/scripts/fetch-pi-rootfs.sh"
}

verify_armv6_elf() {
  local bin="$1"
  if ! command -v readelf >/dev/null; then
    echo "WARNING: readelf not found, skipping ARMv6 check for ${bin}"
    return 0
  fi
  local attrs
  attrs="$(readelf -A "${bin}" 2>/dev/null || true)"
  if echo "${attrs}" | grep -q 'Tag_CPU_arch: v7'; then
    echo "ERROR: ${bin} contains ARMv7 CGO code (Tag_CPU_arch: v7) — not safe on Pi Zero W."
    echo "${attrs}" | grep 'Tag_CPU' || true
    exit 1
  fi
  echo "ARMv6 ELF check passed for ${bin}"
}

run_build() {
  local rootfs="$1"
  local out="$2"
  local goarm="${3:-}"
  local work
  local go_bin
  work="$(prepare_workdir)"
  go_bin="$(chroot_go_bin "${rootfs}")"

  echo "Building ${out} in ${rootfs} (chroot + QEMU, ${go_bin})..."

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

  local -a chroot_env=(
    "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=accept-new"
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

  local host_gopkg
  host_gopkg="$(go env GOPATH)/pkg"
  if mountpoint -q "${rootfs}/root/go/pkg" 2>/dev/null; then
    sudo umount "${rootfs}/root/go/pkg"
  fi

  local go_build_flags=""
  if [ "${goarm}" = "6" ]; then
    sudo rm -rf "${rootfs}/tmp/gocache-goarm6"
    sudo mkdir -p "${rootfs}/root/go/pkg/mod" "${rootfs}/tmp/gocache-goarm6"
    sudo mount --bind "${host_gopkg}/mod" "${rootfs}/root/go/pkg/mod"
    CHROOT_GOMOD="${rootfs}/root/go/pkg/mod"
    chroot_env+=(
      "GOARM=6"
      "GOCACHE=/tmp/gocache-goarm6"
      "PATH=/usr/local/go/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    )
    go_build_flags="-a"
  else
    sudo mkdir -p "${rootfs}/root/go/pkg"
    sudo mount --bind "${host_gopkg}" "${rootfs}/root/go/pkg"
    CHROOT_GOPKG="${rootfs}/root/go/pkg"
    if [ -n "${goarm}" ]; then
      chroot_env+=("GOARM=${goarm}")
    fi
  fi

  mount_chroot_devices "${rootfs}"
  mount_host_resolv "${rootfs}"

  sudo chroot "${rootfs}" env "${chroot_env[@]}" \
    bash -c "cd /src && ${go_bin} build ${go_build_flags} -ldflags '${LDFLAGS}' -o '${out}' cmd/liveserver/main.go"

  cleanup_chroot
  cp "${work}/${out}" "${OUTPUT_DIR}/${out}"
  if [ "${goarm}" = "6" ]; then
    verify_armv6_elf "${OUTPUT_DIR}/${out}"
  fi
}

run_pitest_build() {
  local rootfs="$1"
  local out="$2"
  local goarm="${3:-}"
  local work
  local go_bin
  work="$(prepare_workdir)"
  go_bin="$(chroot_go_bin "${rootfs}")"

  echo "Building ${out} (pitest, CGO disabled) in ${rootfs} (${go_bin})..."

  cleanup_chroot

  local src_mount="${rootfs}/src"
  if mountpoint -q "$src_mount" 2>/dev/null; then
    sudo umount "$src_mount"
  fi
  sudo rm -rf "$src_mount"
  sudo mkdir -p "$src_mount" "${rootfs}/tmp"
  sudo mount --bind "${work}" "$src_mount"
  CHROOT_SRC="$src_mount"

  local -a chroot_env=("CGO_ENABLED=0")
  local go_build_flags=""

  if [ "${goarm}" = "6" ]; then
    sudo rm -rf "${rootfs}/tmp/gocache-pitest-armv6"
    sudo mkdir -p "${rootfs}/tmp/gocache-pitest-armv6"
    chroot_env+=(
      "GOARM=6"
      "GOCACHE=/tmp/gocache-pitest-armv6"
      "PATH=/usr/local/go/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    )
    go_build_flags="-a"
  elif [ -n "${goarm}" ]; then
    chroot_env+=("GOARM=${goarm}")
  fi

  mount_chroot_devices "${rootfs}"
  mount_host_resolv "${rootfs}"

  sudo chroot "${rootfs}" env "${chroot_env[@]}" \
    bash -c "cd /src && ${go_bin} build ${go_build_flags} -o '${out}' cmd/pitest/main.go"

  cleanup_chroot
  cp "${work}/${out}" "${OUTPUT_DIR}/${out}"
}

if [ "$PITEST_ONLY" = 1 ]; then
  if [ ! -f "${BUILD_DIR}/rootfs-armhf/.sgllive-im7-ready" ]; then
    echo "ARM rootfs not found. Running one-time setup..."
    "${ROOT}/scripts/qemu-setup.sh" armhf
  fi
  ensure_raspios_rootfs

  for rootfs in "${BUILD_DIR}/rootfs-armhf" "${BUILD_DIR}/rootfs-raspios"; do
    fix_rootfs_permissions "${rootfs}"
    for mp in "${rootfs}/src" "${rootfs}/tmp/ssh-agent" "${rootfs}/root/go/pkg" "${rootfs}/root/go/pkg/mod" "${rootfs}/etc/resolv.conf" "${rootfs}/dev/null"; do
      if mountpoint -q "$mp" 2>/dev/null; then
        echo "Cleaning stale mount: ${mp}"
        sudo umount "$mp"
      fi
    done
  done

  run_pitest_build "${BUILD_DIR}/rootfs-armhf" pitest_arm32 7
  run_pitest_build "${BUILD_DIR}/rootfs-raspios" pitest_arm32v6 6

  echo "Built:"
  ls -lh "${OUTPUT_DIR}/pitest_arm32" "${OUTPUT_DIR}/pitest_arm32v6"
  exit 0
fi

if [ ! -f "${BUILD_DIR}/rootfs-arm64/.sgllive-im7-ready" ] \
   || [ ! -f "${BUILD_DIR}/rootfs-armhf/.sgllive-im7-ready" ]; then
  echo "Debian ARM rootfs not found. Running one-time setup..."
  "${ROOT}/scripts/qemu-setup.sh" arm64
  "${ROOT}/scripts/qemu-setup.sh" armhf
fi
ensure_raspios_rootfs

for rootfs in "${BUILD_DIR}/rootfs-arm64" "${BUILD_DIR}/rootfs-armhf" "${BUILD_DIR}/rootfs-raspios"; do
  fix_rootfs_permissions "${rootfs}"
  for mp in "${rootfs}/src" "${rootfs}/tmp/ssh-agent" "${rootfs}/root/go/pkg" "${rootfs}/root/go/pkg/mod" "${rootfs}/etc/resolv.conf" "${rootfs}/dev/null"; do
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
run_build "${BUILD_DIR}/rootfs-raspios" liveserver_arm32v6 6

run_pitest_build "${BUILD_DIR}/rootfs-armhf" pitest_arm32 7
run_pitest_build "${BUILD_DIR}/rootfs-raspios" pitest_arm32v6 6

echo "Built:"
ls -lh "${OUTPUT_DIR}/liveserver_arm64" "${OUTPUT_DIR}/liveserver_arm32" "${OUTPUT_DIR}/liveserver_arm32v6" \
       "${OUTPUT_DIR}/pitest_arm32" "${OUTPUT_DIR}/pitest_arm32v6"
