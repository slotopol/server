@echo off
rem This script compiles project for Linux amd64.
rem It produces static C-libraries linkage.
set wd=%~dp0..

xcopy %wd%\confdata %GOPATH%\bin\config /f /d /i /e /k /y

for /F "tokens=*" %%g in ('git describe --tags') do (set buildvers=%%g)
for /F "tokens=*" %%g in ('go run %~dp0\timenow.go') do (set buildtime=%%g)

go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=1
go build -o %GOPATH%\bin\slot_linux_x64 -v -ldflags="-linkmode external -extldflags -static -X 'github.com/slotopol/server/config.BuildVers=%buildvers%' -X 'github.com/slotopol/server/config.BuildTime=%buildtime%'" %wd%
