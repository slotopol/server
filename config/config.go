package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/slotopol/server/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// compiled binary version, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildVers=%buildvers%'"
	BuildVers string

	// compiled binary build date, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildTime=%buildtime%'"
	BuildTime string
)

var (
	// Developer mode, running at debugger.
	DevMode bool
	// AppName is name of this application without extension.
	AppName = util.PathName(os.Args[0])
	// Executable path.
	ExePath string
	// Configuration file with path.
	CfgFile string
	// Configuration path.
	CfgPath string
)

var (
	ErrNoCfgFile = errors.New("configyration file was not found")
)

type CfgJwtAuth struct {
	TokenKey        string        `json:"token-key" yaml:"token-key" mapstructure:"token-key"`
	TokenTimeout    time.Duration `json:"token-timeout" yaml:"token-timeout" mapstructure:"token-timeout"`
	TokenMaxRefresh time.Duration `json:"token-max-refresh" yaml:"token-max-refresh" mapstructure:"token-max-refresh"`
}

type CfgGameplay struct {
	// Maximum value to add to wallet by one transaction.
	AdjunctLimit int
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

func InitConfig() {
	var err error

	if DevMode {
		fmt.Println("*running in developer mode*")
	}
	fmt.Printf("version: %s, builton: %s\n", BuildVers, BuildTime)

	if str, err := os.Executable(); err == nil {
		ExePath = filepath.Dir(str)
	} else {
		ExePath = filepath.Dir(os.Args[0])
	}

	if CfgFile != "" {
		if ok, _ := FileExists(CfgFile); !ok {
			cobra.CheckErr(ErrNoCfgFile)
		}
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		const cfgsub = "config"
		// Search config in home directory with name "slot" (without extension).
		viper.SetConfigName("slot")
		viper.SetConfigType("yaml")
		if env, ok := os.LookupEnv("CFGFILE"); ok {
			viper.AddConfigPath(env)
		}
		viper.AddConfigPath(filepath.Join(ExePath, cfgsub))
		viper.AddConfigPath(ExePath)
		viper.AddConfigPath(cfgsub)
		viper.AddConfigPath(".")
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(filepath.Join(home, cfgsub))
			viper.AddConfigPath(home)
		}
		if env, ok := os.LookupEnv("GOBIN"); ok {
			viper.AddConfigPath(filepath.Join(env, cfgsub))
			viper.AddConfigPath(env)
		} else if env, ok := os.LookupEnv("GOPATH"); ok {
			viper.AddConfigPath(filepath.Join(env, "bin", cfgsub))
			viper.AddConfigPath(filepath.Join(env, "bin"))
		}
	}

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("config file not found!")
	} else {
		cobra.CheckErr(viper.Unmarshal(&Cfg))
		CfgFile = viper.ConfigFileUsed()
		CfgPath = filepath.Dir(CfgFile)
		fmt.Println("config path:", CfgPath)
	}
}

// FileExists check up file existence.
func FileExists(fpath string) (bool, error) {
	var stat, err = os.Stat(fpath)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
