package cmd

import (
	"context"
	"fmt"

	"github.com/slotopol/server/config"
	"github.com/slotopol/server/game/champagne"
	"github.com/slotopol/server/game/dolphinspearl"
	"github.com/slotopol/server/game/jewels"
	"github.com/slotopol/server/game/sizzlinghot"
	"github.com/slotopol/server/game/slotopol"
	"github.com/slotopol/server/game/slotopoldeluxe"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var flags *pflag.FlagSet

const scanShort = "Slots reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for 'Slotopol' game:
  %s scan --slotopol --reels=99.76`

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"calc"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, config.AppName),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx, cancel = context.WithCancel(context.Background())
		defer cancel()

		var is bool
		if is, _ = flags.GetBool("slotopol"); is {
			slotopol.CalcStat(ctx, fReels)
		}
		if is, _ = flags.GetBool("slotopoldeluxe"); is {
			slotopoldeluxe.CalcStat(ctx, fReels)
		}
		if is, _ = flags.GetBool("champagne"); is {
			champagne.CalcStatReg(ctx, fReels)
		}
		if is, _ = flags.GetBool("jewels"); is {
			jewels.CalcStat(ctx, fReels)
		}
		if is, _ = flags.GetBool("sizzlinghot"); is {
			sizzlinghot.CalcStat(ctx, fReels)
		}
		if is, _ = flags.GetBool("dolphinspearl"); is {
			if fReels == "bon" {
				dolphinspearl.CalcStatBon(ctx)
			} else {
				dolphinspearl.CalcStatReg(ctx, fReels)
			}
		}
		return nil
	},
}

var (
	fReels string
)

func init() {
	rootCmd.AddCommand(scanCmd)

	flags = scanCmd.Flags()
	flags.StringVarP(&fReels, "reels", "r", "", "name of reels set to use")

	flags.BoolP("slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	flags.Bool("slotopoldeluxe", false, "'Slotopol Deluxe' Megajack 5x3 slots")
	flags.Bool("champagne", false, "'Champagne' Megajack 5x3 slots")
	flags.Bool("jewels", false, "'Jewels' Novomatic 5x3 slots")
	flags.Bool("sizzlinghot", false, "'Sizzling Hot' Novomatic 5x3 slots")
	flags.Bool("dolphinspearl", false, "'Dolphins Pearl' Novomatic 5x3 slots")
}
