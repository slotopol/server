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

const mcShort = "Slots games reels scanning by Monte Carlo method"
const mcLong = `Calculate RTP (Return to Player) percentage for specified slot game reels by Monte Carlo method.`
const mcExmp = `Scan reels for "Slotopol" game for reels set nearest to 100%%:
  %[1]s mc --game=megajack/slotopol -rtp=100
Scan reels for "Sizzling Hot" game for reels sets nearest to 95.5%% with 5 selected lines and expected precision 0.1%%:
  %[1]s mc -g="Novomatic / Sizzling Hot" -r=95.5 -l=5 --prec=0.1`

var mcflags *pflag.FlagSet

// mcCmd represents the `mc` command
var mcCmd = &cobra.Command{
	Use:     "mc",
	Aliases: []string{"montecarlo"},
	Short:   mcShort,
	Long:    mcLong,
	Example: fmt.Sprintf(mcExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var exitctx = Startup()

		// Load yaml-files
		var noembed bool
		if noembed, err = mcflags.GetBool("noembed"); err != nil {
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
		if alias, err = mcflags.GetString("game"); err != nil {
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
		sp.Method = game.CMmontecarlo
		if sp.TN, err = mcflags.GetInt("mt"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var t uint64
		if t, err = mcflags.GetUint64("total"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		sp.Total = t * 1e6
		var p float64
		if p, err = mcflags.GetFloat64("prec"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		sp.Prec = p / 100
		var c float64
		if c, err = mcflags.GetFloat64("conf"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		sp.Conf = c / 100
		if sp.MRTP, err = mcflags.GetFloat64("rtp"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		if sp.Sel, err = mcflags.GetInt("sel"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		if sp.Sel == 0 {
			sp.Sel = gi.LNum
		} else if sp.Sel > gi.LNum {
			log.Fatalf("number of selected bet lines is greater than maximum number %d in game %s", gi.LNum, alias)
			return
		}
		if sp.Sel != gi.LNum && (gi.GP&game.GPcasc != 0) {
			log.Fatalf("can not change number of selected lines %d on cascade slot %s", gi.LNum, alias)
			return
		}
		scan(exitctx, &sp)
	},
}

func init() {
	rootCmd.AddCommand(mcCmd)

	mcflags = mcCmd.Flags()
	mcflags.Bool("noembed", false, "do not load embedded yaml files, useful for development")
	mcflags.Int("mt", 0, "multithreaded scanning threads number")
	mcflags.StringP("game", "g", "", "identifier of game to scan")
	mcflags.Float64P("rtp", "r", cfg.DefMRTP, "master RTP of game")
	mcflags.IntP("sel", "l", 1, "number of selected bet lines, 0 for all")
	mcflags.Uint64P("total", "n", 1, "Monte Carlo method iterations number, in millions")
	mcflags.Float64("prec", 0.1, "precision of result for Monte Carlo method, in percents")
	mcflags.Float64("conf", 95, "confidence probability, in percents")

	mcCmd.MarkFlagRequired("game")
}
