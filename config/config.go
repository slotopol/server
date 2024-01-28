package config

import (
	"time"
)

var (
	// compiled binary version, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildVers=%buildvers%'"
	BuildVers string

	// compiled binary build date, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildTime=%buildtime%'"
	BuildTime string
)

type CfgJwtAuth struct {
	TokenKey        string        `json:"token-key" yaml:"token-key" mapstructure:"token-key"`
	TokenTimeout    time.Duration `json:"token-timeout" yaml:"token-timeout" mapstructure:"token-timeout"`
	TokenMaxRefresh time.Duration `json:"token-max-refresh" yaml:"token-max-refresh" mapstructure:"token-max-refresh"`
}

type CfgGameplay struct {
	// Maximum value to add to wallet by one transaction.
	AdjunctLimit int `json:"adjunct-limit" yaml:"adjunct-limit" mapstructure:"adjunct-limit"`
	// Maximum number of spin attempts at bad bank balance.
	MaxSpinAttempts int `json:"max-spin-attempts" yaml:"max-spin-attempts" mapstructure:"max-spin-attempts"`
}

type CfgXormDrv struct {
	// Provides XORM driver name.
	XormDriverName string `json:"xorm-driver-name" yaml:"xorm-driver-name" mapstructure:"xorm-driver-name"`
}

// Config is common service settings.
type Config struct {
	CfgJwtAuth  `json:"authentication" yaml:"authentication" mapstructure:"authentication"`
	CfgGameplay `json:"gameplay" yaml:"xorm" mapstructure:"gameplay"`
	CfgXormDrv  `json:"xorm" yaml:"xorm" mapstructure:"xorm"`
}

// Instance of common service settings.
// Inits default values if config is not found.
var Cfg = &Config{
	CfgJwtAuth: CfgJwtAuth{
		TokenKey:        "skJgM4NsbP3fs4k7vh0gfdkgGl8dJTszdLxZ1sQ9ksFnxbgvw2RsGH8xxddUV479",
		TokenTimeout:    1 * 24 * time.Hour,
		TokenMaxRefresh: 3 * 24 * time.Hour,
	},
	CfgGameplay: CfgGameplay{
		AdjunctLimit:    100000,
		MaxSpinAttempts: 300,
	},
	CfgXormDrv: CfgXormDrv{
		XormDriverName: "sqlite3",
	},
}
