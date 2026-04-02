#!/bin/bash

set -euo pipefail

DIR="$(cd "$(dirname "${0}")/../../cmd/kiya" && pwd)"
cd "${DIR}"

goos() {
  case "${1}" in
  Darwin) echo darwin ;;
  Linux) echo linux ;;
  Windows) echo windows ;;
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

LATEST_TAG="$(git tag -l --points-at HEAD | head -n1 || true)"
if [ -z "${LATEST_TAG}" ]; then
  LATEST_TAG="$(git describe --abbrev=0 --tags 2>/dev/null || true)"
fi
VERSION="${LATEST_TAG#v}"
if [ -z "${VERSION}" ]; then
  VERSION="dev"
fi
BUILD_TIME="$(date -u +"%Y%m%dT%H%M%SZ")"
COMMIT="$(git rev-parse --short HEAD)"
FULL_VERSION="${VERSION}+${COMMIT}.${BUILD_TIME}"
echo "Embedded version: $${FULL_VERSION}"

for os in Darwin Linux Windows; do
  for arch in x86_64; do
    dir="${BASE_DIR}/${os}/${arch}/kiya"
    tar_context_dir="$(dirname "${dir}")"
    tar_dir="kiya"
    mkdir -p "${dir}/bin"
    echo GOOS=$(goos "${os}") GOARCH=$(goarch "${arch}") \
      go build \
      -a \
      -ldflags "-X main.version=${FULL_VERSION}" \
      -o "${dir}/bin/kiya" \
      ${SRCS}
    GOOS=$(goos "${os}") GOARCH=$(goarch "${arch}") \
      go build \
      -a \
      -ldflags "-X main.version=${FULL_VERSION}" \
      -o "${dir}/bin/kiya" \
      ${SRCS}
    tar -C "${tar_context_dir}" -cvzf "${BASE_DIR}/kiya-${os}-${arch}.tar.gz" "${tar_dir}"
    cp "${dir}/bin/kiya" "${BASE_DIR}/kiya-${os}-${arch}"
  done
  rm -rf "${BASE_DIR:?/tmp}/${os}"
done
