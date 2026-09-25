#!/bin/bash -u
# This script compiles project for Windows amd64.
# It produces static C-libraries linkage.

wd=$(realpath -s "$(dirname "$0")/..")
mkdir -p "$GOPATH/bin/config" "$GOPATH/bin/sqlite"
cp -ruv "$wd/appdata/"* "$GOPATH/bin/config"

if [[ -z "${BUILDVERS:-}" ]]; then
  BUILDVERS=$(git describe --tags --always)
fi
# See https://tc39.es/ecma262/#sec-date-time-string-format
# time format acceptable for Date constructors.
if [[ -z "${BUILDTIME:-}" ]]; then
  BUILDTIME=$(date +'%FT%T.%3NZ')
fi

go env -w GOOS=windows GOARCH=amd64 CGO_ENABLED=1
go build -o "$GOPATH/bin/slot_win_x64.exe" -v\
 -tags="jsoniter prod full"\
 -buildvcs=false\
 -trimpath -ldflags="-w -s -linkmode external -extldflags -static\
 -X 'github.com/slotopol/server/config.BuildVers=$BUILDVERS'\
 -X 'github.com/slotopol/server/config.BuildTime=$BUILDTIME'"\
 $wd
