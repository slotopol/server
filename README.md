
# Slots server

[![GitHub release][1]][2]
[![Hits-of-Code][3]][4]

[1]: https://img.shields.io/github/v/release/slotopol/server.svg
[2]: https://github.com/slotopol/server/releases/latest
[3]: https://hitsofcode.com/github/slotopol/server?branch=main
[4]: https://hitsofcode.com/github/slotopol/server/view?branch=main

![slotopol-server](docs/logo.webp)

Slots games server. Releases functionality a wide variety of games from different providers like Novomatic, NetEnt, and CT Interactive. Has built-in reels scanner and the sets of reels on different RTP for each game.

Server provides HTTP-based API for popular slots and have well-optimized performance for thousands requests per second. Can be deployed on dedicated server or as portable application for Linux or Windows.

```text
total: 338 games, 170 algorithms, 11 providers
AGT: 64 games
Aristocrat: 6 games
BetSoft: 3 games
CT Interactive: 135 games
IGT: 5 games
Megajack: 4 games
NetEnt: 22 games
Novomatic: 86 games
Play'n GO: 3 games
Playtech: 8 games
Slotopol: 2 games
```

*Last added games*:

* '[Tibetan Song](https://www.livebet2.com/casino/slots/ct-interactive/tibetan-song)' CT Interactive 5x3 videoslot
* '[Black Pharaoh](https://www.livebet2.com/casino/slots/ct-interactive/black-pharaoh)' CT Interactive 5x3 videoslot
* '[Christmas Storm](https://www.livebet2.com/casino/slots/ct-interactive/christmas-storm)' CT Interactive 5x3 cascade videoslot
* '[Beetle Star](https://www.livebet2.com/casino/slots/ct-interactive/beetle-star)' CT Interactive 5x3 videoslot
* '[Amazons Spear](https://www.livebet2.com/casino/slots/ct-interactive/amazons-spear)' CT Interactive 5x3 videoslot
* '[50 Treasures](https://www.livebet2.com/casino/slots/ct-interactive/50-treasures)' CT Interactive 5x4 videoslot
* '[20 Fruitata Wins](https://www.livebet2.com/casino/slots/ct-interactive/20-fruitata-wins)' CT Interactive 5x3 videoslot
* '[20 Star Party](https://www.livebet2.com/casino/slots/ct-interactive/20-star-party)' CT Interactive 5x3 videoslot
* '[30 Fruitata Wins](https://www.livebet2.com/casino/slots/ct-interactive/30-fruitata-wins)' CT Interactive 5x3 cascade videoslot
* '[40 Hell's Cherries](https://www.livebet2.com/casino/slots/ct-interactive/40-hell-s-cherries)' CT Interactive 5x4 videoslot
* '[40 Dice Treasures](https://www.livebet2.com/casino/slots/ct-interactive/40-dice-treasures)' CT Interactive 5x4 videoslot
* '[40 Shining Coins](https://www.livebet2.com/casino/slots/ct-interactive/40-shining-coins)' CT Interactive 5x4 videoslot
* '[40 Shining jewels](https://www.livebet2.com/casino/slots/ct-interactive/40-shining-jewels)' CT Interactive 5x4 videoslot
* '[Celestial Ruler](https://www.livebet2.com/casino/slots/ct-interactive/celestial-ruler)' CT Interactive 5x3 videoslot
* '[Banana Party](https://www.livebet2.com/casino/slots/ct-interactive/banana-party)' CT Interactive 5x3 videoslot
* '[20 Dice Party](https://www.livebet2.com/casino/slots/ct-interactive/20-dice-party)' CT Interactive 5x3 videoslot
* '[Easter Bonanza](https://www.livebet2.com/casino/slots/ct-interactive/easter-bonanza)' CT Interactive 5x4 videoslot
* '[40 Diamond Treasures](https://www.livebet2.com/casino/slots/ct-interactive/40-diamond-treasures)' CT Interactive 5x4 videoslot
* '[Clover Gems](https://www.livebet2.com/casino/slots/ct-interactive/clover-gems)' CT Interactive 5x3 videoslot
* '[Caribbean Adventure](https://www.slotsmate.com/software/ct-interactive/caribbean-adventure)' CT Interactive 5x3 videoslot
* '[Aztec Empress](https://www.slotsmate.com/software/ct-interactive/aztec-empress)' CT Interactive 5x3 videoslot
* '[American Gigolo](https://www.slotsmate.com/software/ct-interactive/american-gigolo)' CT Interactive 5x3 videoslot
* '[Alaska Wild](https://www.slotsmate.com/software/ct-interactive/alaska-wild)' CT Interactive 5x4 videoslot
* '[50 Shining Jewels](https://www.slotsmate.com/software/ct-interactive/50-shining-jewels)' CT Interactive 5x4 videoslot
* '[40 Mega Slot](https://www.slotsmate.com/software/ct-interactive/40-mega-slot)' CT Interactive 5x4 videoslot
* '[40 Roosters](https://www.slotsmate.com/software/ct-interactive/40-roosters)' CT Interactive 5x4 videoslot
* '[30 Treasures](https://www.slotsmate.com/software/ct-interactive/30-treasures)' CT Interactive 5x3 videoslot
* '[Devil's Fruits](https://www.slotsmate.com/software/ct-interactive/devils-fruits)' CT Interactive 3x3 videoslot
* '[Hell's Sevens](https://www.slotsmate.com/software/ct-interactive/hells-sevens)' CT Interactive 3x3 videoslot
* '[Hell's Cherries](https://www.slotsmate.com/software/ct-interactive/hells-cherries)' CT Interactive 3x3 videoslot
* '[Golden Acorn](https://www.slotsmate.com/software/ct-interactive/golden-acorn)' CT Interactive 5x4 videoslot
* '[Roaring Wilds](https://www.slotsmate.com/software/novomatic/roaring-wilds)' Novomatic 5x4 videoslot
* '[Dolphin Reef](https://www.slotsmate.com/software/playtech/dolphin-reef)' Playtech 5x3 videoslot
* '[Rodeo Power](https://www.slotsmate.com/software/ct-interactive/rodeo-power)' CT Interactive 5x3 videoslot
* '[20 Mega Star](https://www.slotsmate.com/software/ct-interactive/20-mega-star)' CT Interactive 5x3 cluster videoslot
* '[Rainbow Charm](https://www.slotsmate.com/software/ct-interactive/rainbow-charm)' CT Interactive 5x3 cluster videoslot
* '[The Magic Goblet](https://www.slotsmate.com/software/ct-interactive/the-magic-goblet)' CT Interactive 5x3 cluster videoslot

*Games with most interesting algorithms released on server*:

* 'Aztec Gold' Megajack 5x3 videoslot has unique 'Aztec Pyramid' bonus game
* '[Book of Ra Deluxe](https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe)' Novomatic 5x3 videoslot with special expanding symbol on free games
* '[Gonzo's Quest](https://www.slotsmate.com/software/netent/gonzos-quest)' NetEnt 5x3 cascade videoslot
* '[Fruit Shop](https://www.slotsmate.com/software/netent/fruit-shop)' NetEnt 5x3 videoslot with free spins on every win
* '[Sun City](https://agtsoftware.com/games/agt/suncity)' AGT 5x3 videoslot with free games continous by the luck
* '[Egypt](https://agtsoftware.com/games/agt/egypt)' AGT 5x3 videoslot with minislot for wilds multiplier
* 'Fire Joker' Play'n GO 5x3 videoslot with jackpot on lined pays and big symbols on free games
* '[Secret Elixir](https://www.slotsmate.com/software/novomatic/secret-elixir)' Novomatic 5x3 videoslot with multipliers on 5th reel
* '[Caribbean Adventure](https://www.slotsmate.com/software/ct-interactive/caribbean-adventure)' CT Interactive 5x3 videoslot, all symbols pays as scatters on freegames
* '[Vikings Fun](https://www.slotsmate.com/software/ct-interactive/vikings-fun)' CT Interactive 5x3 videoslot, some symbols become wilds on free games
* '[Ultra Sevens](https://www.slotsmate.com/software/novomatic/ultra-sevens)' Novomatic 5x4 videoslot with 3 jackpots

## How to build from sources

*Note: you can download the compiled binaries for Windows at [release](https://github.com/slotopol/server/releases/latest) section, or build docker image by [dockerfile](https://github.com/slotopol/server/blob/main/Dockerfile), or use compiled image from [docker hub](https://hub.docker.com/r/schwarzlichtbezirk/slotopol): `docker pull schwarzlichtbezirk/slotopol` and start image [as it described](docs/docker-config.md).*

1. Install [Golang](https://go.dev/dl/) of last version and GCC toolchain *(tip: for Windows it can be MinGW-w64 distribution from [WinLibs](https://winlibs.com/))*.
2. Clone project and download dependencies.
3. Build project with script at `task` directory.

Project can be [built](.github/workflows/go.yml) for Windows/Linux/MacOS.

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

Then web-service can be started:

```cmd
slot_win_x64 -v web
```

The [list of all](docs/list-all.md) provided games can be obtained by command:

```cmd
slot_win_x64 list
```

To find out the list of available RTPs can be with command like this

```cmd
slot_win_x64 list -i megajack --rtp
```

See `slot_win_x64 list -h` with full list of available command line parameters for list-command with filters.

## How to test

To calculate nearest precise RTP to given value for any game on embedded reels, you can use the command like

```cmd
slot_win_x64 scan -g=ctinteractive/luckydollar@95.0
```

Some games algorithms contains Lua-scripts with mathematical RTP calculation without full reels brute force, and returns result on a fly.

You can provide your own reels sets for any game at external yaml-file with same data structure as it in the sources by command

```cmd
slot_win_x64 -f=/some/path/reeldev.yaml web
```

To test gameplay workflow, build [bot](https://github.com/slotopol/bot) as it described, and run some scripts at `script` folder of project. See [readme](https://github.com/slotopol/bot/blob/main/README.md) for details.

## Architecture and logic

**Database.** Service instance oriented on monopoly usage of it's database. It reads necessary database tables on start and avoids any `select` requests at all. Then it stores to database only changes and new data (`update` & `insert`). Those queries are buffered across API endpoints calls to increase performance with database conversations.

It can be used embedded *sqlite* database engine, or *MySQL*, or *PostgreSQL* databases, its configured at `slot-app.yaml` settings file, and by default sqlite is used. Embedded sqlite engine useful for instance started on portable storage, such as flash drive or external SSD, and can serve small sets of players, 50-500 players at the same time. For big number of players it can be used dedicated server with MySQL or PostgreSQL on fast hosted disk.

**Clubs.** There is can be served several clubs. Each club have its own undepended bank, jackpot fund with rate to this fund from spins, and deposit. Bank of club is current balance of club to which arrives coins from users spins, and from which they gets a wins. There is exist linkage of users wins to bank: if bank have not enough coins to pay the win during users spins, this win combination will be skipped. Deposit of club does not used in games, it can be useful to transfer the coins from bank to fix the yield.

**Users accounts.** Accounts have registrations data only. Each account can be associated with several clubs. Each user can have properties for each club with individual balance to gamble, individual access rights, and individual master RTP to choose reels at games.

Each user can play several games at the same time. Each started game have game ID related to user ID and club ID. Any game actions ties to game ID.

## How to use HTTP API

Any API endpoints can receive data in JSON, XML, YAML, or TOML format, depended by `Content-Type` header. If `Content-Type` header not given, JSON will be used to decode as default. `Accept` header if it given, defines response data format. If it absent, same format as at request will be used.

In most cases used `POST`-method of HTTP.

If any response have HTTP status >= 400, body in this case contains error object with `what` field with message and unique source point error `code`.

### Without authorization

First of all you can test replies existence and service is running:

```sh
curl -i localhost:8080/ping
```

Any riplies have `Server` header with content like `slotopol/v0.10.0 (windows; amd64)`. By prefix `slotopol/` you can detect that service is our.

Then you can get a list of games supported by server. This call can be without authorization.

```sh
curl -X GET localhost:8080/game/algs
```

Response has array with available algorithms descriptions. Each structure has a list of games aliases, that shares one algorithm. Field `rtp` has the list of reels with predefined RTP. There is example of structure with info:

```json
,{"aliases":[{"prov":"NetEnt","name":"Trolls","year":2009},{"prov":"NetEnt","name":"Excalibur","year":2013},{"prov":"NetEnt","name":"Pandora's Box","year":2009},{"prov":"NetEnt","name":"Wild Witches","year":2010}],"gt":1,"gp":4628497,"sx":5,"sy":3,"sn":14,"ln":20,"rtp":[87.788791,89.230191,91.925079,93.061471,93.903358,95.183523,96.6485,98.193276,101.929305,110.298257]},
```

To get filtered list of games can be used call like

```sh
curl -X GET localhost:8080/game/list?inc=megajack+ct
```

Where `inc` contains space separated list of filters like in `list` command line parameters, see `slot_win_x64 list -h` for details.

`/servinfo` and `/memusage`, `/signis`, `/sendcode`, `/activate`, `/signup` and `/signin` endpoints also does not expects authorization.

### Authorization

There is supported basic authorization and bearer authorization (with JWT-tokens). Authorization data can be provided by 4 ways: in header `Authorization`, at query parameters, at cookies, and at post form.

**Basic** expects credentials pair `email:password` encoded in unpadded base64 encoding for URL (see RFC 4648).

**Bearer** works with two HS256 JWT tokens - access token and refresh token. Access token should be provided in all cases except `refresh` call. When access-token expires, it should be replaced to refresh-token for refresh-call.

In `/signin` call password can be given by two ways:

1) Explicitly at field `secret` as is.
2) By HMAC SHA256 hash and temporary public key (without send opened secret).

In second case it should be string in field `sigtime` with current time formatted in RFC3339 (can be with nanoseconds). And at field `hs256` it should be hexadecimal HMAC formed with algorithm SHA256 with this current time as a key, and password, i.e. sha256.hmac(sigtime, password). Allowed timeout for public key is 2m 30s.

* Sign-in, and use token from response with any followed calls.

```sh
curl -H "Content-Type: application/json" -d '{"email":"player@example.org","secret":"iVI05M"}' -X POST localhost:8080/signin
```

You can use token `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzbG90b3BvbCIsImV4cCI6NDg2NzQ0NzYxNywibmJmIjoxNzA2NjQ3NjE3LCJ1aWQiOjN9.6g2Hig9ErG8IbvzkPppry5F8HJsMunZPwuQzmetGh4c` for test purpose, it given for user with UID=3 on 100 years. Replace `{{token}}` at samples below to this value.

* When your access token expires (you can get response with 401 status code), use refresh-call with refresh-token to get new tokens pair.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -X GET localhost:8080/refresh
```

### Create and play the game

* Create new game. GID received at response will be used at all calls for access to this game instance. Also you will get initial game state, and user balance at this club.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"cid":1,"uid":3,"alias":"Novomatic/Joker Dolphin"}' -X POST localhost:8080/game/new
```

Endpoint receives `alias` identifier to game that represents as concatenation of provider name and game name with slash. For example, `NetEnt/Tiki Wonders`. Whole list of all supported games can be obtained by [list](docs/list-all.md) command. Identifier turns to lowercase without spaces.

* Change number of selected bet lines.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"gid":1,"sel":5}' -X POST localhost:8080/slot/sel/set
```

* Make a spin. Spin returns `sid` - spin ID, by this ID it can be found at the log; `screen` with new symbols after spin; `wins` with list of win on each line if it was; `fs` - free spins remained; `gain` - total gain after spin, that can be gambled on double up; `wallet` - user balance after spin with won coins.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"gid":1}' -X POST localhost:8080/slot/spin
```

At parameters can be also given new bet value and number of selected lines, if there is no free spins now. So, data at the query can be seen such as `{"gid":1,"bet":2,"sel":5}`. Its can be useful for conversations only by spin-queries.

* Double-up. If presents `gain` after spin, it can be multiplied by gamble. `mult` at argument is multiplier, and it will be `2` for red-black cards game. Returned `gain` will be multiplied on win, and zero on lose. `wallet` represents new balance of user.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"gid":1,"mult":2}' -X POST localhost:8080/slot/doubleup
```

* Collect the gain. After win on spin, or after double-up gain can be collected. In most cases it will be collected automatically on new spin.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"gid":1}' -X POST localhost:8080/slot/collect
```

* Get information about opened game. Recivies GID. Response has game name, game state, last spin ID, user ID, club ID, and user balance at this club.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"gid":1}' -X POST localhost:8080/game/info
```

### Work with user account

* Check up user account existence. It can be done by email or user identifier (`uid` parameter). Call returns true `uid` and `email` if account is found, or zero user identifier if account does not registered.

```sh
curl -X GET localhost:8080/signis?email=rob@example.org
```

* Register new user. E-mail and secret key (password) are expected, name can be omitted. Receives user ID on success.

```sh
curl -H "Content-Type: application/json" -d '{"email":"rob@example.org","secret":"LtpkAr","name":"rob"}' -X POST localhost:8080/signup
```

After registration new user account expects account activation by code sent to user email. Activation can be done by `/activate` endpoint call. If registration of new user was done with admin token, this new user account does not expects activation.

* Activate new user account. It can be done by code sent to user account email. Activation should be done in 15 minutes timeout after registration. If timeout expired activation can be done with new code, sent to email by `/sendcode` endpoint call.

```sh
curl -X GET localhost:8080/activate?uid=3&code=048814
```

Instead `uid` parameter with user identifier can be used user `email`. If activation endpoint was called with admin token, activation code have no matter.

* Send new activation code to user email.

```sh
curl -X GET localhost:8080/sendcode?uid=3
```

Instead `uid` parameter with user identifier can be used user `email`.

* Rename user.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"uid":3,"name":"erigone"}' -X POST localhost:8080/user/rename
```

* Change secret key.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"uid":3,"oldsecret":"iVI05M","newsecret":"pGjKsd"}' -X POST localhost:8080/user/secret
```

* Delete user. Delete-call removes account from database, move all remained user's coins to deposit, and removes all users games from database.

```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer {{token}}" -d '{"uid":3,"secret":"iVI05M"}' -X POST localhost:8080/user/delete
```

---
(c) schwarzlichtbezirk, 2024-2025.
