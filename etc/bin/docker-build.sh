#!/usr/bin/env bash
set -eu

PLATFORM="${PLATFORM:-linux/amd64}"
IMAGE="${IMAGE:-github-kramphub-kiya}"

LATEST_TAG="$(git tag -l --points-at HEAD | head -n1 || true)"
if [ -z "$LATEST_TAG" ]; then
  LATEST_TAG="$(git describe --abbrev=0 --tags 2>/dev/null || true)"
fi

VERSION="${LATEST_TAG#v}"
if [ -z "$VERSION" ]; then
  VERSION="dev"
fi

BUILD_TIME="$(date -u +"%Y%m%dT%H%M%SZ")"
COMMIT="$(git rev-parse --short HEAD)"
FULL_VERSION="${VERSION}+${COMMIT}.${BUILD_TIME}"

echo "Building image $IMAGE:$VERSION"
echo "Embedded version: $FULL_VERSION"

docker buildx build \
  --no-cache \
  --platform "$PLATFORM" \
  --build-arg VERSION="$FULL_VERSION" \
  -t "$IMAGE:$VERSION" \
  --progress=plain \
  --load .
