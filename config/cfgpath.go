package cfg

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/slotopol/server/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Prints more verbose information to log.
	Verbose int
	// Silent mode, turn off all log output.
	Silent bool
	// AppName is name of this application without extension.
	AppName = util.PathName(os.Args[0])
	// Executable path.
	ExePath string
	// Configuration file with path.
	CfgFile string
	// Configuration path.
	CfgPath string
	// SQLite-files path.
	SqlPath string
	// List of paths to yaml files with game specific data to load on start up.
	ObjPath []string
)

var (
	ErrNoCfgFile = errors.New("configuration file was not found")
)

func InitConfig() {
	var err error

	InitLogger()

	Infof("version: %s, builton: %s", BuildVers, BuildTime)

	ExePath = func() string {
		if str, err := os.Executable(); err == nil {
			return filepath.Dir(str)
		} else {
			return filepath.Dir(os.Args[0])
		}
	}()

	// Config path
	if val := os.Getenv("SLOTOPOL_CFGFILE"); val != "" {
		CfgFile, _ = util.ExpandHomePath(val)
	}
	if CfgFile != "" {
		if ok, _ := FileExists(CfgFile); !ok {
			cobra.CheckErr(ErrNoCfgFile)
		}
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		const sub = "appdata"
		// Search config in home directory with name "slot-config" (without extension).
		viper.SetConfigName("slot-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(filepath.Join(ExePath, sub))
		viper.AddConfigPath(ExePath)
		viper.AddConfigPath(sub)
		viper.AddConfigPath(".")
		if appdata, err := os.UserConfigDir(); err == nil {
			viper.AddConfigPath(filepath.Join(appdata, "slotopol", "server"))
		}
		if env, ok := os.LookupEnv("GOBIN"); ok {
			viper.AddConfigPath(filepath.Join(env, sub))
			viper.AddConfigPath(env)
		} else if env, ok := os.LookupEnv("GOPATH"); ok {
			viper.AddConfigPath(filepath.Join(env, "bin", sub))
			viper.AddConfigPath(filepath.Join(env, "bin"))
		}
	}

	viper.BindEnv("authentication.access-key", "SLOTOPOL_ACCESSKEY")
	viper.BindEnv("authentication.refresh-key", "SLOTOPOL_REFRESHKEY")
	viper.BindEnv("activation.brevo-api-key", "SLOTOPOL_BREVOAPIKEY")
	viper.BindEnv("database.driver-name", "SLOTOPOL_DBDRIVER")
	viper.BindEnv("database.club-source-name", "SLOTOPOL_CLUBDSN")
	viper.BindEnv("database.spin-source-name", "SLOTOPOL_SPINDSN")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		Warnf("config file not found!")
	} else {
		cobra.CheckErr(viper.Unmarshal(&Cfg))
		CfgFile = viper.ConfigFileUsed()
		CfgPath = filepath.Dir(CfgFile)
		Infof("config path: %s", CfgPath)
	}

	// Check up signing keys existence and length.
	if len(Cfg.AccessKey) == 0 {
		cobra.CheckErr(errors.New("access-key for JWT-token signing does not given, please set SLOTOPOL_ACCESSKEY environment variable or set access-key in config file"))
	}
	if len(util.S2B(Cfg.AccessKey)) < 32 {
		cobra.CheckErr(errors.New("access-key for JWT-token signing is too short, it should be at least 32 bytes long"))
	}
	if len(Cfg.RefreshKey) == 0 {
		cobra.CheckErr(errors.New("refresh-key for JWT-token signing does not given, please set SLOTOPOL_REFRESHKEY environment variable or set refresh-key in config file"))
	}
	if len(util.S2B(Cfg.RefreshKey)) < 32 {
		cobra.CheckErr(errors.New("refresh-key for JWT-token signing is too short, it should be at least 32 bytes long"))
	}
	if Cfg.AccessKey == Cfg.RefreshKey {
		cobra.CheckErr(errors.New("access-key can not be equal to refresh-key"))
	}
	if Cfg.UseActivation {
		if len(Cfg.BrevoApiKey) == 0 {
			cobra.CheckErr(errors.New("Brevo API-key for email sending does not given, please set SLOTOPOL_BREVOAPIKEY environment variable or set brevo-api-key in config file"))
		}
	}

	// Detect SQLite path.
	if val := os.Getenv("SLOTOPOL_SQLPATH"); val != "" {
		SqlPath, _ = util.ExpandHomePath(val)
		cobra.CheckErr(os.MkdirAll(SqlPath, os.ModePerm))
	}
	if SqlPath == "" {
		SqlPath = CfgPath
	}
	Infof("sqlite path: %s", SqlPath)
}

// DirExists check up directory existence.
func DirExists(fpath string) (bool, error) {
	var stat, err = os.Stat(fpath)
	if err == nil {
		return stat.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
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
