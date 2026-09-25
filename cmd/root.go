package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

const rootShort = "Slots games backend"
const rootLong = `This application implements web server and reels scanner for slots games.`

var FinalPaths []string

var (
	rootCmd = &cobra.Command{
		Use:     cfg.AppName,
		Version: cfg.BuildVers,
		Short:   rootShort,
		Long:    rootLong,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			for _, pattern := range cfg.ObjPath {
				var pattern = filepath.ToSlash(pattern)
				pattern, err := util.ExpandHomePath(pattern)
				if err != nil {
					return fmt.Errorf("failed to expand home path for %q: %w", pattern, err)
				}
				pattern = util.Envfmt(pattern, nil)
				matches, err := doublestar.FilepathGlob(pattern,
					doublestar.WithCaseInsensitive(),
					doublestar.WithFilesOnly(),
					doublestar.WithNoHidden(),
				)
				if err != nil {
					if errors.Is(err, doublestar.ErrPatternNotExist) {
						return fmt.Errorf("path or its base directory does not exist: %q", pattern)
					}
					return fmt.Errorf("invalid pattern %q: %w", pattern, err)
				}
				FinalPaths = append(FinalPaths, matches...)
				cfg.Infof("found %d files for pattern %q", len(matches), pattern)
			}
			return nil
		}}
)

func init() {
	cobra.OnInitialize(cfg.InitConfig)

	var pf = rootCmd.PersistentFlags()
	pf.StringVarP(&cfg.CfgFile, "config", "c", "", "config file (default is config/slot-app.yaml at executable location)")
	pf.StringVarP(&cfg.SqlPath, "sqlite", "q", "", "sqlite databases path (default same as config file path)")
	pf.StringSliceVarP(&cfg.ObjPath, "fpath", "f", nil, "additional paths/patterns to yaml files or folders with game specific data")
	pf.CountVarP(&cfg.Verbose, "verbose", "v", "print more verbose information to log, can repeated to increase verbosity level, for example: -vvv")
	pf.BoolVar(&cfg.Silent, "silent", false, "turn off all log output")
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s, builton: %s\n", cfg.BuildVers, cfg.BuildTime))
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
