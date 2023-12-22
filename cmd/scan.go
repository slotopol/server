package cmd

import (
	"github.com/schwarzlichtbezirk/slot-srv/game/slotopol"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Slots reels scanning",
	Long:  `Calculate RTP (Return to Player) percentage for specified slot game reels.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if fSlotopol {
			slotopol.CalcStat()
		}
		return nil
	},
}

var (
	fSlotopol bool
)

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolVarP(&fSlotopol, "slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
}
