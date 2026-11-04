#!/bin/bash -u
# This script compiles project with Megajack games only.
# It produces static C-libraries linkage.

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

goos=$(go env GOOS)
goarch=$(go env GOARCH)
if [[ "$goarch" == "amd64" ]]; then
  goarch="x64"
fi
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  appext=".exe"
else
  appext=""
fi
go env -w CGO_ENABLED=1
go build -o "${GOPATH}/bin/slot_${goos}_${goarch}_megajack${appext}" -v \
  -tags="jsoniter prod megajack" \
  -buildvcs=false \
  -trimpath -ldflags="-w -s -linkmode external -extldflags -static \
  -X 'github.com/slotopol/server/config.BuildVers=$BUILDVERS' \
  -X 'github.com/slotopol/server/config.BuildTime=$BUILDTIME'" \
  $wd
