package cmd

import (
	"context"
	"fmt"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/config/links"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var flags *pflag.FlagSet

const scanShort = "Slots reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for 'Slotopol' game:
  %s scan --slotopol --reels=100`

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"calc"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, cfg.AppName),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var exitctx context.Context
		if exitctx, err = Init(); err != nil {
			return
		}

		for _, iter := range links.ScatIters {
			iter(flags, exitctx)
		}

		return
	},
}

var (
	fReels string
)

func init() {
	rootCmd.AddCommand(scanCmd)

	flags = scanCmd.Flags()
	flags.StringVarP(&fReels, "reels", "r", "", "name of reels set to use")

	for _, setter := range links.FlagsSetters {
		setter(flags)
	}
}
