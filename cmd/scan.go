package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const scanShort = "Slots games reels scanning"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels.`
const scanExmp = `Scan reels for "Slotopol" game for reels set nearest to 100%%:
  %[1]s scan --game=megajack/slotopol@100
Scan reels for "Dolphins Pearl" and "Katana" games for reels sets nearest to 92%% and 94.5%% respectively:
  %[1]s scan -g="Novomatic / Dolphins Pearl @ 92" -g=novomatic/katana@94.5`

var scanflags *pflag.FlagSet

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

		var list []string
		if list, err = scanflags.GetStringArray("game"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var lstage bool
		if lstage, err = scanflags.GetBool("lstage"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var vstage bool
		if vstage, err = scanflags.GetBool("vstage"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var gi *game.GameInfo
		var scan game.Scanner
		var ok bool
		for i, rid := range list {
			var rs = strings.Split(rid, "@")
			var alias = strings.TrimSpace(rs[0])
			var mrtp float64
			if len(rs) > 1 {
				if mrtp, err = strconv.ParseFloat(strings.TrimSpace(rs[1]), 64); err != nil {
					log.Fatalf("can not parse master RTP for '%s': %s", alias, err.Error())
					return
				}
			} else {
				mrtp = cfg.DefMRTP
			}
			var aid = util.ToID(alias)
			if gi, ok = game.InfoMap[aid]; !ok {
				log.Fatalf("game name \"%s\" does not recognized", alias)
				return
			}
			if len(gi.RTP) == 0 {
				log.Fatalf("RTP list does not complete for %s", alias)
				return
			}
			if scan, ok = game.ScanFactory[aid]; !ok {
				log.Fatalf("game name \"%s\" does not recognized", alias)
				return
			}
			if scan == nil {
				fmt.Println()
				fmt.Printf("*** scanner for '%s' game does not provided ***\n", alias)
			}
			if len(list) > 1 {
				fmt.Println()
				var msg = fmt.Sprintf("*** (%d/%d) scan '%s' game with master RTP %g ***", i+1, len(list), alias, mrtp)
				if lstage {
					log.Println(msg)
				}
				if vstage {
					fmt.Println(msg)
				}
			}
			scan(exitctx, mrtp)
			if exitctx.Err() != nil {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanflags = scanCmd.Flags()
	scanflags.StringArrayP("game", "g", nil, "identifier of game to scan")
	scanflags.Bool("noembed", false, "do not load embedded yaml files, useful for development")
	scanflags.Bool("lstage", false, "log verbose stage information during scanning")
	scanflags.Bool("vstage", true, "print verbose stage information during scanning")
	scanflags.Uint64Var(&cfg.MCCount, "mc", 0, "Monte Carlo method samples number, in millions")
	scanflags.Float64Var(&cfg.MCPrec, "mcp", 0, "Precision of result for Monte Carlo method, in percents")
	scanflags.IntVar(&cfg.MTCount, "mt", 0, "multithreaded scanning threads number")

	scanCmd.MarkFlagRequired("game")
}
