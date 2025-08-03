#!/bin/bash -u
# This script compiles project for Windows x86.
# It produces static C-libraries linkage.

wd=$(realpath -s "$(dirname "$0")/..")
mkdir -p "$GOPATH/bin/config" "$GOPATH/bin/sqlite"
cp -ruv "$wd/appdata/"* "$GOPATH/bin/config"

buildvers=$(git describe --tags)
# See https://tc39.es/ecma262/#sec-date-time-string-format
# time format acceptable for Date constructors.
buildtime=$(date +'%FT%T.%3NZ')

go env -w GOOS=windows GOARCH=386 CGO_ENABLED=1
go build -o "$GOPATH/bin/slot_win_x86.exe" -v\
 -tags="jsoniter prod full"\
 -buildvcs=false\
 -ldflags="-linkmode external -extldflags -static\
 -X 'github.com/slotopol/server/config.BuildVers=$buildvers'\
 -X 'github.com/slotopol/server/config.BuildTime=$buildtime'"\
 $wd
