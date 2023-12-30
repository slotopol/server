package cmd

import (
	"context"
	"fmt"

	"github.com/slotopol/server/config"
	"github.com/slotopol/server/game/champagne"
	"github.com/slotopol/server/game/dolphinspearl"
	"github.com/slotopol/server/game/sizzlinghot"
	"github.com/slotopol/server/game/slotopol"
	"github.com/slotopol/server/game/slotopoldeluxe"
	"github.com/spf13/cobra"
)

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

		if fSlotopol {
			slotopol.CalcStat(ctx, fReels)
		}
		if fSlotopolDeluxe {
			slotopoldeluxe.CalcStat(ctx, fReels)
		}
		if fChampagne {
			champagne.CalcStatReg(ctx, fReels)
		}
		if fSizzlingHot {
			sizzlinghot.CalcStat(ctx, fReels)
		}
		if fDolphinsPearl {
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
	fReels          string
	fSlotopol       bool
	fSlotopolDeluxe bool
	fChampagne      bool
	fSizzlingHot    bool
	fDolphinsPearl  bool
)

func init() {
	rootCmd.AddCommand(scanCmd)

	var flags = scanCmd.Flags()
	flags.StringVarP(&fReels, "reels", "r", "", "name of reels set to use")
	flags.BoolVarP(&fSlotopol, "slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	flags.BoolVar(&fSlotopolDeluxe, "slotopoldeluxe", false, "'Slotopol Deluxe' Megajack 5x3 slots")
	flags.BoolVar(&fChampagne, "champagne", false, "'Champagne' Megajack 5x3 slots")
	flags.BoolVar(&fSizzlingHot, "sizzlinghot", false, "'Sizzling Hot' Novomatic 5x3 slots")
	flags.BoolVar(&fDolphinsPearl, "dolphinspearl", false, "'Dolphins Pearl' Novomatic 5x3 slots")
}
