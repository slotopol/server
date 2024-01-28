package cmd

import (
	"fmt"

	"github.com/slotopol/server/config"
	"github.com/spf13/cobra"
)

const rootShort = "Slots games backend"
const rootLong = `This application performs all implemented tasks for slots games.`

var (
	rootCmd = &cobra.Command{
		Use:     config.AppName,
		Version: config.BuildVers,
		Short:   rootShort,
		Long:    rootLong,
	}
)

func init() {
	cobra.OnInitialize(config.InitConfig)

	var flags = rootCmd.PersistentFlags()
	flags.StringVarP(&config.CfgFile, "config", "c", "", "config file (default is config/slot.yaml at executable location)")
	flags.StringVarP(&config.SqlPath, "sqlite", "q", "", "sqlite databases path (default same as config file path)")
	flags.BoolVarP(&config.DevMode, "devmode", "d", false, "start application in developer mode")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", config.BuildVers, config.BuildTime))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
