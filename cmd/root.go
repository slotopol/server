package cmd

import (
	"fmt"

	cfg "github.com/slotopol/server/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const rootShort = "Slots games backend"
const rootLong = `This application implements web server and reels scanner for slots games.`

var rootflags *pflag.FlagSet

var (
	rootCmd = &cobra.Command{
		Use:     cfg.AppName,
		Version: cfg.BuildVers,
		Short:   rootShort,
		Long:    rootLong,
	}
)

func init() {
	cobra.OnInitialize(cfg.InitConfig)

	rootflags = rootCmd.PersistentFlags()
	rootflags.StringVarP(&cfg.CfgFile, "config", "c", "", "config file (default is config/slot-app.yaml at executable location)")
	rootflags.StringVarP(&cfg.SqlPath, "sqlite", "q", "", "sqlite databases path (default same as config file path)")
	rootflags.StringArrayVarP(&cfg.ObjPath, "fpath", "f", nil, "additional paths to yaml files or folders with game specific data (can be repeated)")
	rootflags.BoolVarP(&cfg.DevMode, "devmode", "d", false, "start application in developer mode")
	rootflags.BoolVarP(&cfg.Verbose, "verbose", "v", false, "print more verbose information to log")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", cfg.BuildVers, cfg.BuildTime))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
