package cmd

import (
	"fmt"
	"os"

	"github.com/schwarzlichtbezirk/slot-srv/config"
	"github.com/spf13/cobra"
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

var (
	rootCmd = &cobra.Command{
		Use:     BaseName(os.Args[0]),
		Version: config.BuildVers,
		Short:   "Slots games backend",
		Long:    `This application performs all implemented tasks for slots games.`,
	}
)

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s", config.BuildVers, config.BuildTime))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
