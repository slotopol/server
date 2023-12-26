package cmd

import (
	"fmt"

	"github.com/schwarzlichtbezirk/slot-srv/config"
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
			slotopol.CalcStat()
		}
		if fSlotopolDeluxe {
			slotopoldeluxe.CalcStat()
		}
		return nil
	},
}

var (
	fSlotopol       bool
	fSlotopolDeluxe bool
)

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolVarP(&fSlotopol, "slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	scanCmd.Flags().BoolVar(&fSlotopolDeluxe, "slotopoldeluxe", false, "'Slotopol Deluxe' Megajack 5x3 slots")
}
