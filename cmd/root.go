package cmd

import (
	"fmt"

	"github.com/schwarzlichtbezirk/slot-srv/config"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     config.AppName,
		Version: config.BuildVers,
		Short:   "Slots games backend",
		Long:    `This application performs all implemented tasks for slots games.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("***")
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "config file (default is config/slot.yaml at executable location)")
	rootCmd.PersistentFlags().BoolVarP(&config.DevMode, "devmode", "d", false, "start application in developer mode")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", config.BuildVers, config.BuildTime))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
