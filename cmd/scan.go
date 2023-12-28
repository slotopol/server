package cmd

import (
	"fmt"

	"github.com/schwarzlichtbezirk/slot-srv/config"
	"github.com/schwarzlichtbezirk/slot-srv/game/dolphinspearl"
	"github.com/schwarzlichtbezirk/slot-srv/game/sizzlinghot"
	"github.com/schwarzlichtbezirk/slot-srv/game/slotopol"
	"github.com/schwarzlichtbezirk/slot-srv/game/slotopoldeluxe"
	"github.com/spf13/cobra"
)

const scanShort = "Slots reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for 'Slotopol' game:
  %s scan -s`

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"calc"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, config.AppName),
	RunE: func(cmd *cobra.Command, args []string) error {
		if fSlotopol {
			slotopol.CalcStat(fReels)
		}
		if fSlotopolDeluxe {
			slotopoldeluxe.CalcStat(fReels)
		}
		if fSizzlingHot {
			sizzlinghot.CalcStat(fReels)
		}
		if fDolphinsPearl {
			if fReels == "bon" {
				dolphinspearl.CalcStatBon()
			} else {
				dolphinspearl.CalcStatReg(fReels)
			}
		}
		return nil
	},
}

var (
	fReels          string
	fSlotopol       bool
	fSlotopolDeluxe bool
	fSizzlingHot    bool
	fDolphinsPearl  bool
)

func init() {
	rootCmd.AddCommand(scanCmd)

	var flags = scanCmd.Flags()
	flags.StringVarP(&fReels, "reels", "r", "", "name of reels set to use")
	flags.BoolVarP(&fSlotopol, "slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	flags.BoolVar(&fSlotopolDeluxe, "slotopoldeluxe", false, "'Slotopol Deluxe' Megajack 5x3 slots")
	flags.BoolVar(&fSizzlingHot, "sizzlinghot", false, "'Sizzling Hot' Novomatic 5x3 slots")
	flags.BoolVar(&fDolphinsPearl, "dolphinspearl", false, "'Dolphins Pearl' Novomatic 5x3 slots")
}
