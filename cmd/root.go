package cmd

import (
	"fmt"

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

	var flags = rootCmd.PersistentFlags()
	flags.StringVarP(&cfg.CfgFile, "config", "c", "", "config file (default is config/slot-app.yaml at executable location)")
	flags.StringVarP(&cfg.SqlPath, "sqlite", "q", "", "sqlite databases path (default same as config file path)")
	flags.BoolVarP(&cfg.DevMode, "devmode", "d", false, "start application in developer mode")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", cfg.BuildVers, cfg.BuildTime))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
