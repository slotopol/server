package cmd

import (
	"fmt"

	"github.com/slotopol/server/config"
	game "github.com/slotopol/server/game"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var scanflags *pflag.FlagSet

const scanShort = "Slots games reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for 'Slotopol' game for reels set nearest to 100%%:
  %[1]s scan --slotopol --reels=100`

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"calc"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, config.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var exitctx = Startup()

		for _, iter := range game.ScanIters {
			iter(scanflags, exitctx)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanflags = scanCmd.Flags()
	scanflags.Float64P("reels", "r", 92.5, "master RTP to calculate nearest reels")
	scanflags.Uint64Var(&config.MCCount, "mc", 0, "Monte Carlo method samples number, in millions")
	scanflags.BoolVar(&config.MTScan, "mt", false, "multithreaded scanning")

	for _, gi := range game.GameList {
		for _, ga := range gi.Aliases {
			scanflags.Bool(ga.ID, false, fmt.Sprintf("'%s' %s %dx%d videoslot", ga.Name, gi.Provider, gi.SX, gi.SY))
		}
	}
	for _, setter := range game.FlagsSetters {
		setter(scanflags)
	}
}
