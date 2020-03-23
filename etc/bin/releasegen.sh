#!/bin/bash

set -euo pipefail

DIR="$(cd "$(dirname "${0}")/../../cmd/kiya" && pwd)"
cd "${DIR}"

goos() {
  case "${1}" in
  Darwin) echo darwin ;;
  Linux) echo linux ;;
  *) return 1 ;;
  esac
}

goarch() {
  case "${1}" in
  x86_64) echo amd64 ;;
  *) return 1 ;;
  esac
}

BASE_DIR="../../release"
rm -rf "${BASE_DIR}"

SRCS=$(find . -type f -name "*.go" -maxdepth 1 | grep -v "test")

for os in Darwin Linux Windows; do
  for arch in x86_64; do
    dir="${BASE_DIR}/${os}/${arch}/kiya"
    tar_context_dir="$(dirname "${dir}")"
    tar_dir="kiya"
    mkdir -p "${dir}/bin"
    CGO_ENABLED=0 GOOS=$(goos "${os}") GOARCH=$(goarch "${arch}") \
      go build \
      -a \
      -installsuffix cgo \
      -ldflags "-X 'main.version=$(git tag -l --points-at HEAD)'" \
      -o "${dir}/bin/kiya" \
      ${SRCS}
    tar -C "${tar_context_dir}" -cvzf "${BASE_DIR}/kiya-${os}-${arch}.tar.gz" "${tar_dir}"
    cp "${dir}/bin/kiya" "${BASE_DIR}/kiya-${os}-${arch}"
  done
  rm -rf "${BASE_DIR:?/tmp}/${os}"
done
