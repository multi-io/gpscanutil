#!/usr/bin/env bash
# Usage: render-formula.sh VERSION CHECKSUMS_FILE [DOWNLOAD_URL]
# Prints the Homebrew formula for the release archives listed in
# CHECKSUMS_FILE (goreleaser's checksums.txt).
set -euo pipefail

version=$1
checksums=$2
url=${3:-https://github.com/multi-io/gpscanutil/releases/download/v$version}

args=(-e "s|@VERSION@|$version|g" -e "s|@URL@|$url|g")
for platform in darwin_arm64 darwin_amd64 linux_arm64 linux_amd64; do
  sha=$(awk -v f="gpscanutil_${version}_${platform}.tar.gz" '$2 == f { print $1 }' "$checksums")
  if [[ -z $sha ]]; then
    echo "no checksum for $platform in $checksums" >&2
    exit 1
  fi
  args+=(-e "s|@SHA256_${platform}@|$sha|g")
done

sed "${args[@]}" "$(dirname "$0")/gpscanutil.rb.in"
