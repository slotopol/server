package cmd

import (
	"fmt"
	"os"

	cfg "github.com/slotopol/server/config"

	"github.com/spf13/cobra"
)

const rootShort = "Slots games backend"
const rootLong = `This application implements web server and reels scanner for slots games.`

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

	var pf = rootCmd.PersistentFlags()
	pf.StringVarP(&cfg.CfgFile, "config", "c", "", "config file (default is config/slot-app.yaml at executable location)")
	pf.StringVarP(&cfg.SqlPath, "sqlite", "q", "", "sqlite databases path (default same as config file path)")
	pf.StringArrayVarP(&cfg.ObjPath, "fpath", "f", nil, "additional paths to yaml files or folders with game specific data (can be repeated)")
	pf.BoolVarP(&cfg.Verbose, "verbose", "v", false, "print more verbose information to log")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", cfg.BuildVers, cfg.BuildTime))
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
