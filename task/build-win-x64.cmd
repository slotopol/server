@echo off
rem This script compiles project for Windows amd64.
rem It produces static C-libraries linkage.

set wd=%~dp0..
xcopy %wd%\appdata %GOPATH%\bin\config /f /d /i /e /k /y

for /F "tokens=*" %%g in ('git describe --tags') do (set buildvers=%%g)
for /f "tokens=2 delims==" %%g in ('wmic os get localdatetime /value') do set dt=%%g
set buildtime=%dt:~0,4%-%dt:~4,2%-%dt:~6,2%T%dt:~8,2%:%dt:~10,2%:%dt:~12,2%.%dt:~15,3%Z

go env -w GOOS=windows GOARCH=amd64 CGO_ENABLED=1
go build -o %GOPATH%\bin\slot_win_x64.exe -v -tags="jsoniter prod full" -buildvcs=false -ldflags="-linkmode external -extldflags -static -X 'github.com/slotopol/server/config.BuildVers=%buildvers%' -X 'github.com/slotopol/server/config.BuildTime=%buildtime%'" %wd%
