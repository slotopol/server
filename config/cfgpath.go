package cfg

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/slotopol/server/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var (
	// Prints more verbose information to log.
	Verbose bool
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
	// Multithreaded scanning threads number
	MTCount int
	// Monte Carlo method samples number, in millions
	MCCount uint64
	// Precision of result for Monte Carlo method, in percents
	MCPrec float64
)

var (
	XormStorage *xorm.Engine
	XormSpinlog *xorm.Engine
)

var (
	ErrNoCfgFile = errors.New("configuration file was not found")
)

func InitConfig() {
	var err error

	if Verbose {
		log.Printf("version: %s, builton: %s\n", BuildVers, BuildTime)
	}

	ExePath = func() string {
		if str, err := os.Executable(); err == nil {
			return filepath.Dir(str)
		} else {
			return filepath.Dir(os.Args[0])
		}
	}()

	// Config path
	if val := os.Getenv("SLOTOPOL_CFGFILE"); val != "" {
		CfgFile = val
	}
	if CfgFile != "" {
		if ok, _ := FileExists(CfgFile); !ok {
			cobra.CheckErr(ErrNoCfgFile)
		}
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		const sub = "config"
		// Search config in home directory with name "slot-app" (without extension).
		viper.SetConfigName("slot-app")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(filepath.Join(ExePath, sub))
		viper.AddConfigPath(ExePath)
		viper.AddConfigPath(sub)
		viper.AddConfigPath("appdata")
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

	viper.BindEnv("database.driver-name", "SLOTOPOL_DBDRIVER")
	viper.BindEnv("database.club-source-name", "SLOTOPOL_CLUBDSN")
	viper.BindEnv("database.spin-source-name", "SLOTOPOL_SPINDSN")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Println("config file not found!")
	} else {
		cobra.CheckErr(viper.Unmarshal(&Cfg))
		CfgFile = viper.ConfigFileUsed()
		CfgPath = filepath.Dir(CfgFile)
		if Verbose {
			log.Printf("config path: %s\n", CfgPath)
		}
	}

	// Detect SQLite path.
	if SqlPath == "" {
		SqlPath = LookupInLocations("SLOTOPOL_SQLPATH", "sqlite", "slot-club.sqlite")
	}
	if Verbose {
		log.Printf("sqlite path: %s\n", SqlPath)
	}
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

func LookupInLocations(env, sub, fname string) (fpath string) {
	var list []string
	if val := os.Getenv(env); val != "" {
		fpath = val
		cobra.CheckErr(os.MkdirAll(fpath, os.ModePerm))
		return
	}
	list, _ = AddDir(list,
		filepath.Join(ExePath, sub),
		ExePath,
		sub,
		"appdata",
		".",
	)
	if appdata, err := os.UserCacheDir(); err == nil {
		list, _ = AddDir(list, filepath.Join(appdata, "slotopol", "server"))
	}
	if env, ok := os.LookupEnv("GOBIN"); ok {
		list, _ = AddDir(list, filepath.Join(env, sub))
		list, _ = AddDir(list, env)
	} else if env, ok := os.LookupEnv("GOPATH"); ok {
		list, _ = AddDir(list, filepath.Join(env, "bin", sub))
		list, _ = AddDir(list, filepath.Join(env, "bin"))
	}
	if fpath = LookupDir(list, fname); fpath == "" {
		fpath = filepath.Join(ExePath, sub)
		cobra.CheckErr(os.MkdirAll(fpath, os.ModePerm))
	}
	return
}

func LookupDir(list []string, fname string) string {
	for _, fpath := range list {
		if ok, _ := FileExists(filepath.Join(fpath, fname)); ok {
			return fpath
		}
	}
	return ""
}

func AbsDir(dir string) (string, error) {
	dir = os.ExpandEnv(dir)
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir), nil
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return dir, err
	}
	return filepath.Clean(dir), nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AddDir(list []string, dirs ...string) ([]string, error) {
	var errs []error
	var err error
	for _, dir := range dirs {
		if dir, err = AbsDir(dir); err != nil {
			errs = append(errs, err)
			continue
		}
		if StringInSlice(dir, list) {
			continue
		}
		list = append(list, dir)
	}
	return list, errors.Join(errs...)
}
