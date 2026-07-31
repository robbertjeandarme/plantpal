#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="${root}/dist"

mkdir -p "${out}"

build() {
  local goos=$1 goarch=$2 name=$3
  echo "→ ${name}"
  CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
    go build -ldflags="-s -w" -o "${out}/${name}" .
}

cd "${root}"

build linux   amd64 plantpal-linux-amd64
build linux   arm64 plantpal-linux-arm64
build darwin  amd64 plantpal-darwin-amd64
build darwin  arm64 plantpal-darwin-arm64
build windows amd64 plantpal-windows-amd64.exe
build windows arm64 plantpal-windows-arm64.exe

echo
echo "Binaries in ${out}:"
ls -lh "${out}"
