#!/bin/bash -u
# This script compiles project for Linux amd64 inside of docker.
# It produces static C-libraries linkage.

# dockerfile has no access to git repository,
# so update content of this variable by
#   echo $(git describe --tags)
buildvers="v0.7.0-6-gb2085ae"
# See https://tc39.es/ecma262/#sec-date-time-string-format
# time format acceptable for Date constructors.
buildtime=$(date +'%FT%T.%3NZ')

go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=1
go build -o /go/bin/app -v\
 -tags="jsoniter prod full"\
 -ldflags="-linkmode external -extldflags -static\
 -X 'github.com/slotopol/server/config.BuildVers=$buildvers'\
 -X 'github.com/slotopol/server/config.BuildTime=$buildtime'"\
 ./
