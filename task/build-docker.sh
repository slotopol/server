#!/bin/bash -u
# This script compiles project for Linux amd64 inside of docker.
# It produces static C-libraries linkage.

wd=$(realpath -s "$(dirname "$0")/..")
mkdir -p "$GOPATH/bin/appdata"
cp -ruv "$wd/appdata/"* "$GOPATH/bin/appdata"

# dockerfile has no access to git repository,
# so set up BUILDVERS outside of dockerfile call.
#   docker build --build-arg "BUILDVERS=$(git describe --tags --always)" -t schwarzlichtbezirk/slotopol:latest .
if [[ -z "${BUILDVERS:-}" ]]; then
  BUILDVERS=$(git describe --tags --always) # try to get build version from git if not set
  if [[ -z "${BUILDVERS:-}" ]]; then
    echo "Warning: BUILDVERS is not set and git describe failed. Use --build-arg 'BUILDVERS=$(git describe --tags --always)'"
    BUILDVERS="unknown"
  fi
fi
# See https://tc39.es/ecma262/#sec-date-time-string-format
# time format acceptable for Date constructors.
if [[ -z "${BUILDTIME:-}" ]]; then
  BUILDTIME=$(date +'%FT%T.%3NZ')
fi

go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=1
go build -o /go/bin/app -v \
  -tags="jsoniter prod full" \
  -buildvcs=false \
  -trimpath -ldflags="-w -s -linkmode external -extldflags -static \
  -X 'github.com/slotopol/server/config.BuildVers=$BUILDVERS' \
  -X 'github.com/slotopol/server/config.BuildTime=$BUILDTIME'" \
  $wd
