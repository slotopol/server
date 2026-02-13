package cmd

import (
	"fmt"
	"log"
	"runtime"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const scanShort = "Slots games reels scanning by brute force"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels by brute force method.`
const scanExmp = `Scan reels for "Slotopol" game for reels set nearest to 100%%:
  %[1]s scan --game=megajack/slotopol -rtp=100
Scan reels for "Dolphins Pearl" game for reels sets nearest to 92%% with 10 selected lines:
  %[1]s scan -g="Novomatic / Dolphins Pearl" -r=92 -l=10`

var scanflags *pflag.FlagSet

// scanCmd represents the `scan` command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"bruteforce"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var exitctx = Startup()

		// Load yaml-files
		var noembed bool
		if noembed, err = scanflags.GetBool("noembed"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		if !noembed {
			LoadInternalYaml(exitctx)
		}
		if err = LoadExternalYaml(exitctx); err != nil {
			log.Fatalf("can not load yaml files: %s", err.Error())
			return
		}
		UpdateAlgList()

		var alias string
		if alias, err = scanflags.GetString("game"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var aid = util.ToID(alias)
		var gi *game.GameInfo
		var ok bool
		if gi, ok = game.InfoMap[aid]; !ok {
			log.Fatalf("game name \"%s\" does not recognized", alias)
			return
		}
		if len(gi.RTP) == 0 {
			log.Fatalf("RTP list does not complete for %s", alias)
			return
		}

		var scan game.Scanner
		if scan, ok = game.ScanFactory[aid]; !ok {
			log.Fatalf("game name \"%s\" does not recognized", alias)
			return
		}
		if scan == nil {
			fmt.Println()
			fmt.Printf("*** scanner for '%s' game does not provided ***\n", alias)
		}

		var sp game.ScanPar
		sp.Method = game.CMbruteforce

		var tn int
		if tn, err = scanflags.GetInt("mt"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		if tn < 1 {
			tn = runtime.GOMAXPROCS(0)
		}
		sp.TN = tn

		if sp.MRTP, err = scanflags.GetFloat64("rtp"); err != nil {
			log.Fatalln(err.Error())
			return
		}

		if gi.LNum > 0 {
			var sel int
			if sel, err = scanflags.GetInt("sel"); err != nil {
				log.Fatalln(err.Error())
				return
			}
			if sel == 0 {
				sel = gi.LNum
			} else if sel > gi.LNum {
				log.Fatalf("number of selected bet lines is greater than maximum number %d in game %s", gi.LNum, alias)
				return
			}
			if sel != gi.LNum && (gi.GP&game.GPcasc != 0) {
				log.Fatalf("can not change number of selected lines %d on cascade slot %s", gi.LNum, alias)
				return
			}
			sp.Sel = sel
		}

		scan(exitctx, &sp)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanflags = scanCmd.Flags()
	scanflags.Bool("noembed", false, "do not load embedded yaml files, useful for development")
	scanflags.Int("mt", 0, "multithreaded scanning threads number")
	scanflags.StringP("game", "g", "", "identifier of game to scan")
	scanflags.Float64P("rtp", "r", cfg.DefMRTP, "master RTP of game")
	scanflags.IntP("sel", "l", 1, "number of selected bet lines, 0 for all")

	scanCmd.MarkFlagRequired("game")
}
