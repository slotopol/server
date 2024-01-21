
# Slot server

[![Go](https://github.com/slotopol/server/actions/workflows/go.yml/badge.svg)](https://github.com/slotopol/server/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/v/release/slotopol/server.svg)](https://github.com/slotopol/server/releases/latest)
[![Hits-of-Code](https://hitsofcode.com/github/slotopol/server?branch=main)](https://hitsofcode.com/github/slotopol/server/view?branch=main)

Slots games server.

This project is at an early stage of development.

# How to build

1. Install [Golang](https://go.dev/dl/) of last version.
2. Clone project and download dependencies.
3. Build project with script at `task` directory.

For Windows command prompt:

```cmd
git clone https://github.com/slotopol/server.git
cd server
go mod download && go mod verify
task\build-win-x64.cmd
```

or for Linux shell or git bash:

```sh
git clone https://github.com/slotopol/server.git
cd server
go mod download && go mod verify
sudo chmod +x ./task/*.sh
./task/build-linux-x64.sh
```
