package config

var (
	// compiled binary version, sets by compiler with command
	//    go build -ldflags="-X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildVers=%buildvers%'"
	BuildVers string

	// compiled binary build date, sets by compiler with command
	//    go build -ldflags="-X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildTime=%buildtime%'"
	BuildTime string
)
