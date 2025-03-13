
# Run modes of docker image

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

In second way volume with sqlite files can be pointed by environment variable `SQLPATH`:

```sh
docker run -d -p 8080:8080 -v "D:/srv":/cache -e SQLPATH=/cache/sqlite schwarzlichtbezirk/slotopol
```

Default settings file can also be replaced. To replace only one settings file can by command:

```sh
docker run -d -p 8080:8080 -v "D:/srv/config/slot-app.yaml":/go/bin/config/slot-app.yaml:ro schwarzlichtbezirk/slotopol
```

If you mount folder volume with settings files, other configuration files also should be present there. In case of mounted folder with configuration files you can keep several settings files with different names in one folder, and provide specified file by `CFGFILE` environment variable. I.e. you can make `slot-app-mysql.yaml` file with connection to MySQL, and `slot-app-postgres.yaml` file with PostgreSQL connection for example.

```sh
docker run -d -p 8080:8080 -v //d/srv:/cache -e CFGFILE=/cache/config/slot-app.yaml -e SQLPATH=/cache/sqlite schwarzlichtbezirk/slotopol
```
