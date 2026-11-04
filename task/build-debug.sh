#!/bin/bash -u
# This script compiles project for default platform with debug information.
# It produces static C-libraries linkage.

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}
need go
need git

wd=$(realpath -s "$(dirname "$0")/..")
mkdir -p "$GOPATH/bin/appdata"
cp -ruv "$wd/appdata/"* "$GOPATH/bin/appdata"

if [[ -z "${BUILDVERS:-}" ]]; then
  BUILDVERS=$(git describe --tags --always)
fi
# See https://tc39.es/ecma262/#sec-date-time-string-format
# time format acceptable for Date constructors.
if [[ -z "${BUILDTIME:-}" ]]; then
  BUILDTIME=$(date +'%FT%T.%3NZ')
fi

go env -w CGO_ENABLED=1
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  appext=".exe"
else
  appext=""
fi
go build -o "${GOPATH}/bin/slot_debug${appext}" -v \
  -tags="jsoniter" \
  -buildvcs=false \
  -ldflags="-linkmode external -extldflags -static \
  -X 'github.com/slotopol/server/config.BuildVers=$BUILDVERS' \
  -X 'github.com/slotopol/server/config.BuildTime=$BUILDTIME'" \
  $wd
