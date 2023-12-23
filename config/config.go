package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// compiled binary version, sets by compiler with command
	//    go build -ldflags="-X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildVers=%buildvers%'"
	BuildVers string

	// compiled binary build date, sets by compiler with command
	//    go build -ldflags="-X 'github.com/schwarzlichtbezirk/slot-srv/config.BuildTime=%buildtime%'"
	BuildTime string
)

var (
	// Executable path.
	ExePath string
	// Configuration file with path.
	CfgFile string
	// Configuration path.
	CfgPath string
	// AppName is name of this application without extension.
	AppName = BaseName(os.Args[0])
	// Developer mode, running at debugger.
	DevMode bool
)

// BaseName returns name of file in given file path without extension.
func BaseName(fpath string) string {
	var j = len(fpath)
	if j == 0 {
		return ""
	}
	var i = j - 1
	for {
		if os.IsPathSeparator(fpath[i]) {
			i++
			break
		}
		if fpath[i] == '.' {
			j = i
		}
		if i == 0 {
			break
		}
		i--
	}
	return fpath[i:j]
}

func InitConfig() {
	var err error

	if DevMode {
		fmt.Println("*running in developer mode*")
	}

	if str, err := os.Executable(); err == nil {
		ExePath = path.Dir(str)
	} else {
		ExePath = path.Dir(os.Args[0])
	}

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		const cfgsub = "config"
		// Search config in home directory with name "slot" (without extension).
		viper.SetConfigName("slot")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path.Join(ExePath, cfgsub))
		viper.AddConfigPath(ExePath)
		viper.AddConfigPath(cfgsub)
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/" + cfgsub)
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("$GOPATH/bin/" + cfgsub)
		viper.AddConfigPath("$GOPATH/bin")
	}

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err == nil {
		CfgFile = viper.ConfigFileUsed()
		fmt.Println("Using config file:", CfgFile)
		CfgPath = filepath.Dir(CfgFile)

		if err = viper.Unmarshal(&Cfg); err != nil {
			cobra.CheckErr(err)
		}
	}
}

// Config is common service settings.
type Config struct {
}

// Instance of common service settings.
var Cfg = &Config{ // inits default values:
}
