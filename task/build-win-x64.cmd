@echo off
rem This script compiles project for Windows amd64.
rem It produces static C-libraries linkage.
set wd=%~dp0..

xcopy %wd%\confdata %GOPATH%\bin\config /f /d /i /e /k /y

for /F "tokens=*" %%g in ('git describe --tags') do (set buildvers=%%g)
for /F "tokens=*" %%g in ('go run %~dp0\timenow.go') do (set buildtime=%%g)

go env -w GOOS=windows GOARCH=amd64 CGO_ENABLED=1
go build -o %GOPATH%\bin\slot_win_x64.exe -v -ldflags="-linkmode external -extldflags -static -X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildVers=%buildvers%' -X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildTime=%buildtime%'" %wd%
