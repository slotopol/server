package cmd

import (
	"fmt"
	"log"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var scanflags *pflag.FlagSet

const scanShort = "Slots games reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for "Slotopol" game for reels set nearest to 100%%:
  %[1]s scan --game=megajack/slotopol --mrtp=100
Scan reels for "Dolphins Pearl" and "Katana" games for reels set nearest to 94.5%%:
  %[1]s scan -g="Novomatic / Dolphins Pearl" -g=novomatic/katana -r=94.5`

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"calc"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var exitctx = Startup()

		var mrtp float64
		if mrtp, err = scanflags.GetFloat64("mrtp"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var list []string
		if list, err = scanflags.GetStringArray("game"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var scan game.Scanner
		var ok bool
		for _, alias := range list {
			if scan, ok = game.ScanFactory[util.ToID(alias)]; !ok {
				log.Fatalf("game name \"%s\" does not recognized", alias)
				return
			}
			if scan == nil {
				fmt.Println()
				fmt.Printf("***Scanner for '%s' game is absent***\n", alias)
			}
			if len(list) > 1 {
				fmt.Println()
				fmt.Printf("***Scan '%s' game with master RTP %g***\n", alias, mrtp)
			}
			scan(exitctx, mrtp)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanflags = scanCmd.Flags()
	scanflags.StringArrayP("game", "g", nil, "identifier of game to scan")
	scanflags.Float64P("mrtp", "r", cfg.DefMRTP, "master RTP to calculate nearest reels")
	scanflags.Uint64Var(&cfg.MCCount, "mc", 0, "Monte Carlo method samples number, in millions")
	scanflags.Float64Var(&cfg.MCPrec, "mcp", 0, "Precision of result for Monte Carlo method, in percents")
	scanflags.IntVar(&cfg.MTCount, "mt", 0, "multithreaded scanning threads number")

	scanCmd.MarkFlagRequired("game")
}
