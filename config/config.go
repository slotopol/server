package config

import (
	"errors"
	"fmt"
	"os"
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
	// Developer mode, running at debugger.
	DevMode bool
	// AppName is name of this application without extension.
	AppName = PathName(os.Args[0])
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

// Config is common service settings.
type Config struct {
}

// Instance of common service settings.
// Inits default values if config is not found.
var Cfg = &Config{}

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

// PathName returns name of file in given file path without extension.
func PathName(fpath string) string {
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
