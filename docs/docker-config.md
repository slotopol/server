
# Run modes of docker image

Image can be used to start stand-alone container with embedded sqlite database, or start container connected to external database. External database can be `MySQL` or `PostgreSQL`, and can be located in docker network or outside.

## Start stand-alone container

In shortest way image can be started by command

```sh
docker run -d -p 8080:8080 schwarzlichtbezirk/slotopol
```

By default it uses embedded `sqlite3` database engine, and writes database files inside of container. Configuration files with default settings also placed inside of container.

To keep database files outside of container can by two ways. For example you have some directory on host machine where should be its placed: `D:/srv/sqlite`.

In first way can be mounted volume with this path and replace default path:

```sh
docker run -d -p 8080:8080 -v "D:/srv/sqlite":/go/bin/sqlite schwarzlichtbezirk/slotopol
```

or in bash-like syntax:

```sh
docker run -d -p 8080:8080 -v //d/srv/sqlite:/go/bin/sqlite schwarzlichtbezirk/slotopol
```

In second way volume with sqlite files can be pointed by environment variable `SLOTOPOL_SQLPATH`:

```sh
docker run -d -p 8080:8080 -v "D:/srv":/cache -e SLOTOPOL_SQLPATH=/cache/sqlite schwarzlichtbezirk/slotopol
```

Default settings file can also be replaced. To replace only one settings file can by command:

```sh
docker run -d -p 8080:8080 -v "D:/srv/config/slot-app.yaml":/go/bin/config/slot-app.yaml:ro schwarzlichtbezirk/slotopol
```

If you mount folder volume with settings files, other configuration files also should be present there. In case of mounted folder with configuration files you can keep several settings files with different names in one folder, and provide specified file by `SLOTOPOL_CFGFILE` environment variable. I.e. you can make `slot-app-mysql.yaml` file with connection to MySQL, and `slot-app-postgres.yaml` file with PostgreSQL connection for example.

```sh
docker run -d -p 8080:8080 -v //d/srv:/cache -e SLOTOPOL_CFGFILE=/cache/config/slot-app.yaml -e SLOTOPOL_SQLPATH=/cache/sqlite schwarzlichtbezirk/slotopol
```

## Start container with database connection

There is environment variables can be used to configure database connection to MySQL or PostgreSQL:

- `SLOTOPOL_DBDRIVER`

Database driver name, overwrites `database.driver-name` config setting in yaml-file. It can be "sqlite3", "mysql", "postgres". To place database in memory only, point driver name "sqlite3" and source names ":memory:".

- `SLOTOPOL_CLUBDSN`

Overwrites `database.club-source-name` config setting. Data source name for 'club' database to create XORM engine. For "sqlite3" driver it should be db file name without path like `slot-club.sqlite`. For "mysql" driver it should match to pattern of standard MySQL DSN: `user:password@tcp(addr:port)/slot_club`. For postgres it should match to pattern of PostgreSQL connection string: `host=localhost port=5432 user=admin password=secret dbname=slot_spin sslmode=disable`.

- `SLOTOPOL_SPINDSN`

Overwrites `database.spin-source-name` config setting. Data source name for 'spin' database to create XORM engine. Values could fit to same pattern as at previous case. Main mean of this setting is to split sensitive data and logs to splitted databases (`slot_club` and `slot_spin` for example), and put data at previous setting and logs at this setting. If Slotopol service connected to database engine in docker container it can be useful both settings refer to one database.

Logs can be switched off by `SLOTOPOL_SPINDSN` is set to empty string.

## Sample of docker-compose with MySQL connection in docker

```yaml
# Basic use case to run Slotopol server with MySQL database in docker.
# Typical usage:
#   docker compose -p slotopol-mysql -f "docker-compose.mysql.yaml" up -d --build

services:

  slotopol:
    image: schwarzlichtbezirk/slotopol:latest
    container_name: slotopol-srv
    depends_on:
      - database
    ports:
      - 8080:8080
    networks:
      slotopol-net:
        aliases:
          - slotopolhost
    environment:
      SLOTOPOL_DBDRIVER: mysql
      SLOTOPOL_CLUBDSN: admin:jB44t78dHv@tcp(mysqlhost:3306)/slotopol
      SLOTOPOL_SPINDSN: admin:jB44t78dHv@tcp(mysqlhost:3306)/slotopol
    restart: unless-stopped
    stop_signal: SIGINT
    stop_grace_period: 15s

  database:
    image: mysql:8.0-debian
    container_name: slotopol-db
    ports:
      - 3306:3306
    networks:
      slotopol-net:
        aliases:
          - mysqlhost
    environment:
      MYSQL_ROOT_PASSWORD: dg8FxCm1hh
      MYSQL_USER: admin
      MYSQL_PASSWORD: jB44t78dHv
      MYSQL_DATABASE: slotopol

networks:
  slotopol-net:
    name: slotopol-net
    external: false
    ipam:
      driver: default
```
