#!/bin/bash
# Build minimal pitest binaries only (no CGO, no AppBackend deps).

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${1:-${ROOT}/liveserver}"

exec "${ROOT}/scripts/build_liveserver.sh" --pitest-only "$OUTPUT_DIR"
