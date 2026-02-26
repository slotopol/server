package cmd

import (
	"fmt"
	"log"
	"runtime"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const scanShort = "Slots games reels scanning by brute force or by Monte Carlo"
const scanLong = `Calculate RTP (Return to Player) percentage for specified slot game reels by brute force method ("scan" or "bf" command) or by Monte Carlo method ("mc" command).`
const scanExmp = `Scan reels for "Slotopol" game for reels set nearest to 100%%:
  %[1]s scan --game=megajack/slotopol -rtp=100
Scan reels for "Dolphins Pearl" game for reels sets nearest to 92%% with 10 selected lines:
  %[1]s bf -g="Novomatic / Dolphins Pearl" -r=92 -l=10
Scan reels with using Monte Carlo method for "Sizzling Hot" game for reels sets nearest to 95.5%% with 5 selected lines and expected precision 0.1%%:
  %[1]s mc -g="Novomatic / Sizzling Hot" -r=95.5 -l=5 --prec=0.1`

func SetupParSheet(pf *pflag.FlagSet, sp *game.ScanPar, gi *game.GameInfo) (err error) {
	var ok bool

	var tn int
	if tn, err = pf.GetInt("mt"); err != nil {
		return
	}
	if tn < 1 {
		tn = runtime.GOMAXPROCS(0)
	}
	sp.TN = tn

	if sp.MRTP, err = pf.GetFloat64("rtp"); err != nil {
		return
	}

	if gi.LNum > 0 {
		var sel int
		if sel, err = pf.GetInt("sel"); err != nil {
			return
		}
		if sel == 0 {
			sel = gi.LNum
		} else if sel > gi.LNum {
			return fmt.Errorf("number of selected bet lines is greater than maximum number %d in game %s", gi.LNum, gi.ID())
		}
		if sel != gi.LNum && (gi.GP&game.GPcasc != 0) {
			return fmt.Errorf("can not change number of selected lines %d on cascade slot %s", gi.LNum, gi.ID())
		}
		sp.Sel = sel
	}

	var c float64
	if c, err = pf.GetFloat64("conf"); err != nil {
		return
	}
	sp.Conf = c / 100

	var t uint64
	if t, err = pf.GetUint64("total"); err != nil {
		return
	}
	sp.Total = t * 1e6

	var p float64
	if p, err = pf.GetFloat64("prec"); err != nil {
		return
	}
	sp.Prec = p / 100

	var pfm = map[string]uint{
		"main":    slot.PF_main,
		"fg":      slot.PF_fg,
		"vi":      slot.PF_vi,
		"ci":      slot.PF_ci,
		"ranges":  slot.PF_ranges,
		"contrib": slot.PF_contrib,
		"plain":   slot.PF_plain,
	}
	for sf, uf := range pfm {
		if ok, err = pf.GetBool(sf); err != nil {
			return
		}
		if ok {
			sp.PF |= uf
		}
	}
	return
}

// scanCmd represents the `scan` command
var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"bf", "mc", "bruteforce", "montecarlo"},
	Short:   scanShort,
	Long:    scanLong,
	Example: fmt.Sprintf(scanExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var exitctx = Startup()
		var pf = cmd.Flags()

		// Load yaml-files
		var noembed bool
		if noembed, err = pf.GetBool("noembed"); err != nil {
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
		if alias, err = pf.GetString("game"); err != nil {
			log.Fatalln(err.Error())
			return
		}
		var aid = util.ToID(alias)
		var gi *game.GameInfo
		var ok bool
		if gi, ok = game.InfoMap[aid]; !ok {
			log.Fatalf("game name \"%s\" does not recognized", gi.ID())
			return
		}
		if len(gi.RTP) == 0 {
			log.Fatalf("RTP list does not complete for %s", gi.ID())
			return
		}

		var scan game.Scanner
		if scan, ok = game.ScanFactory[aid]; !ok {
			log.Fatalf("game name \"%s\" does not recognized", gi.ID())
			return
		}
		if scan == nil {
			fmt.Println()
			fmt.Printf("*** scanner for '%s' game does not provided ***\n", gi.ID())
		}

		var sp game.ScanPar

		switch cmd.CalledAs() {
		case "scan", "bf", "bruteforce":
			sp.Method = game.CMbruteforce
		case "mc", "montecarlo":
			sp.Method = game.CMmontecarlo
		}

		if err = SetupParSheet(pf, &sp, gi); err != nil {
			log.Fatalln(err.Error())
			return
		}

		scan(exitctx, &sp)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	var pf = scanCmd.Flags()
	pf.Bool("noembed", false, "do not load embedded yaml files, useful for development")
	pf.StringP("game", "g", "", "identifier of game to scan")
	// ParSheet
	pf.Int("mt", 0, "multithreaded scanning threads number")
	pf.Float64P("rtp", "r", cfg.DefMRTP, "master RTP of game")
	pf.IntP("sel", "l", 0, "number of selected bet lines, 0 for all")
	pf.Float64("conf", 95, "confidence probability, in percents")
	pf.Uint64P("total", "n", 10, "Monte Carlo method iterations number, in millions")
	pf.Float64("prec", 0.1, "precision of result for Monte Carlo method, in percents")
	// print flags
	pf.Bool("main", true, "print RTP, sigma and other main information")
	pf.Bool("fg", true, "print info for bonus reels")
	pf.Bool("vi", true, "print volatility index")
	pf.Bool("ci", true, "print index of convergence")
	pf.Bool("ranges", false, "print RTP ranges")
	pf.Bool("contrib", false, "print symbols contribution to payouts")
	pf.Bool("plain", false, "simulator plain data")

	pf.SortFlags = false

	scanCmd.MarkFlagRequired("game")
}
