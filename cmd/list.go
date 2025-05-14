package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const listShort = "List of available games released on server"
const listLong = ``
const listExmp = `Get the list of all available games:
  %[1]s list
Get the list of slots with cascade falls:
  %[1]s list -i casc
Get the list of 'NetExt' and 'BetSoft' games:
  %[1]s list -i netent -i betsoft
Get the list of Megajack games and any games with 3 reels:
  %[1]s list -i megajack -i 3x
Get the list of games with screen 3x3 without 'AGT' games:
  %[1]s list -i 3x3 -e agt
Get the list of 'Play'n GO' games with RTP list for each:
  %[1]s list -i playngo --rtp`

var listflags *pflag.FlagSet

var (
	fSort            bool
	fProp, fRTP      bool
	fMrtp, fDiff     float64
	inclist, exclist []string
)

type filter func(*game.GameInfo) bool

var FiltMap = map[string]filter{
	"all":  func(gi *game.GameInfo) bool { return true },
	"slot": func(gi *game.GameInfo) bool { return gi.GT == game.GTslot },

	"keno":       func(gi *game.GameInfo) bool { return gi.GT == game.GTkeno },
	"agt":        func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "agt" },
	"aristocrat": func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "aristocrat" },
	"betsoft":    func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "betsoft" },
	"igt":        func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "igt" },
	"megajack":   func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "megajack" },
	"netent":     func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "netent" },
	"novomatic":  func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "novomatic" },
	"playngo":    func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "playngo" },
	"playtech":   func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "playtech" },
	"slotopol":   func(gi *game.GameInfo) bool { return util.ToID(gi.Prov) == "slotopol" },

	"lines": func(gi *game.GameInfo) bool { return gi.LN > 0 },
	"ways":  func(gi *game.GameInfo) bool { return gi.WN > 0 },
	"casc":  func(gi *game.GameInfo) bool { return gi.GT == game.GTslot && gi.GP&game.GPcasc > 0 },
	"jack":  func(gi *game.GameInfo) bool { return gi.GP&game.GPjack > 0 },
	"fg":    func(gi *game.GameInfo) bool { return gi.GP&(game.GPfghas+game.GPretrig) > 0 },
	"bon":   func(gi *game.GameInfo) bool { return gi.BN > 0 },
}
var (
	finclist, fexclist []filter
	reReel             = regexp.MustCompile(`^(\d)x$`)
	reScrn             = regexp.MustCompile(`^(\d)x(\d)$`)
	reYEQ              = regexp.MustCompile(`^y=(\d{2}|\d{4})$`)
	reYLT              = regexp.MustCompile(`^y<(\d{2}|\d{4})$`)
	reYGT              = regexp.MustCompile(`^y>(\d{2}|\d{4})$`)
	reLNEQ             = regexp.MustCompile(`^ln=(\d{1,3})$`)
	reLNLT             = regexp.MustCompile(`^ln<(\d{1,3})$`)
	reLNGT             = regexp.MustCompile(`^ln>(\d{1,3})$`)
	reWNEQ             = regexp.MustCompile(`^wn=(\d{1,4})$`)
	reWNLT             = regexp.MustCompile(`^wn<(\d{1,4})$`)
	reWNGT             = regexp.MustCompile(`^wn>(\d{1,4})$`)
)

func getfilter(key string) filter {
	key = util.ToLower(key)
	if f, ok := FiltMap[key]; ok {
		return f
	}
	if s := reReel.FindStringSubmatch(key); len(s) > 0 {
		var x, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.SX == x }
	}
	if s := reScrn.FindStringSubmatch(key); len(s) > 0 {
		var x, _ = strconv.Atoi(s[1])
		var y, _ = strconv.Atoi(s[2])
		return func(gi *game.GameInfo) bool { return gi.SX == x && gi.SY == y }
	}
	if s := reYEQ.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[1])
		if year < 100 {
			year += 2000
		}
		return func(gi *game.GameInfo) bool { return gi.Year == year }
	}
	if s := reYLT.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[1])
		if year < 100 {
			year += 2000
		}
		return func(gi *game.GameInfo) bool { return gi.Year < year }
	}
	if s := reYGT.FindStringSubmatch(key); len(s) > 0 {
		var year, _ = strconv.Atoi(s[1])
		if year < 100 {
			year += 2000
		}
		return func(gi *game.GameInfo) bool { return gi.Year > year }
	}
	if s := reLNEQ.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.LN == ln }
	}
	if s := reLNLT.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.LN > 0 && gi.LN < ln }
	}
	if s := reLNGT.FindStringSubmatch(key); len(s) > 0 {
		var ln, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.LN > ln }
	}
	if s := reWNEQ.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.WN == wn }
	}
	if s := reWNLT.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.WN > 0 && gi.WN < wn }
	}
	if s := reWNGT.FindStringSubmatch(key); len(s) > 0 {
		var wn, _ = strconv.Atoi(s[1])
		return func(gi *game.GameInfo) bool { return gi.WN > wn }
	}
	return nil
}

func passes(gi *game.GameInfo) bool {
	var is bool
	for _, f := range finclist {
		if f(gi) {
			is = true
			break
		}
	}
	if !is {
		return false
	}
	for _, f := range fexclist {
		if f(gi) {
			return false
		}
	}
	return true
}

func FormatGameInfo(gi *game.GameInfo) string {
	var b strings.Builder
	if gi.SN > 0 {
		if gi.GP&game.GPcasc > 0 {
			fmt.Fprintf(&b, "'%s' %s %dx%d cascade videoslot", gi.Name, gi.Prov, gi.SX, gi.SY)
		} else {
			fmt.Fprintf(&b, "'%s' %s %dx%d videoslot", gi.Name, gi.Prov, gi.SX, gi.SY)
		}
	} else {
		fmt.Fprintf(&b, "'%s' %s %d spots lottery", gi.Name, gi.Prov, gi.SX)
	}
	if fProp {
		if gi.SN > 0 {
			fmt.Fprintf(&b, ", %d symbols", gi.SN)
		}
		if gi.LN > 0 {
			if gi.GP&game.GPlsel == 0 {
				fmt.Fprintf(&b, ", constant %d lines", gi.LN)
			} else {
				fmt.Fprintf(&b, ", %d lines", gi.LN)
			}
		}
		if gi.WN > 0 {
			fmt.Fprintf(&b, ", %d ways", gi.WN)
		}
		if gi.GP&game.GPjack > 0 {
			b.WriteString(", has jackpot")
		}
		if gi.GP&game.GPfill > 0 {
			b.WriteString(", has multiplier on filled screen")
		}
		if gi.GP&(game.GPfghas+game.GPretrig) > 0 {
			b.WriteString(", ")
			if gi.GP&game.GPretrig > 0 {
				b.WriteString("retriggerable ")
			}
			b.WriteString("free games")
			if gi.GP&game.GPfgmult > 0 {
				b.WriteString(" with multiplier")
			}
			if gi.GP&game.GPfgreel > 0 {
				b.WriteString(" on bonus reels")
			}
		}
		if gi.GP&game.GPscat > 0 {
			b.WriteString(", has scatters")
		}
		if gi.GP&game.GPwild > 0 {
			if gi.GP&game.GPwmult > 0 {
				b.WriteString(", has wilds with multiplier")
			} else {
				b.WriteString(", has wilds")
			}
		}
		if gi.GP&game.GPrwild > 0 {
			if gi.GP&game.GPwmult > 0 {
				b.WriteString(", has reel wilds with multiplier")
			} else {
				b.WriteString(", has reel wilds")
			}
		}
		if gi.GP&game.GPbwild > 0 {
			b.WriteString(", has big wilds")
		}
		if gi.GP&game.GPwturn > 0 {
			b.WriteString(", symbols turns to wilds")
		}
		if gi.GP&game.GPbsym > 0 {
			b.WriteString(", has big symbols")
		}
	}
	if len(gi.RTP) > 0 {
		if fMrtp > 0 {
			var rtp = gi.FindClosest(fMrtp)
			fmt.Fprintf(&b, ", mrtp(%g)=%g", fMrtp, rtp)
		}
		if fDiff > 0 {
			var diff = gi.FindClosest(fDiff) - fDiff
			fmt.Fprintf(&b, ", diff(%g)=%.6g", fDiff, diff)
		}
		if fRTP {
			b.WriteString(", RTP: ")
			for i, rtp := range gi.RTP {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(&b, "%.2f", rtp)
			}
		}
	}
	return b.String()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Long:    listLong,
	Example: fmt.Sprintf(listExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var f filter
		for _, key := range inclist {
			if f = getfilter(key); f == nil {
				fmt.Printf("filter with name '%s' does not recognized\n", key)
				continue
			}
			finclist = append(finclist, f)
		}
		for _, key := range exclist {
			if f = getfilter(key); f == nil {
				fmt.Printf("filter with name '%s' does not recognized\n", key)
				continue
			}
			fexclist = append(fexclist, f)
		}

		var alg = map[*game.AlgDescr]int{}
		var prov = map[string]int{}
		var gamelist = make([]*game.GameInfo, 0, 256)
		for _, gi := range game.InfoMap {
			if passes(gi) {
				alg[gi.AlgDescr]++
				prov[gi.Prov]++
				gamelist = append(gamelist, gi)
			}
		}

		if is, _ := listflags.GetBool("name"); is {
			fmt.Println()
			sort.Slice(gamelist, func(i, j int) bool {
				var gii, gij = gamelist[i], gamelist[j]
				if fSort {
					if gii.Prov == gij.Prov {
						return gii.Name < gij.Name
					}
					return gii.Prov < gij.Prov
				} else {
					if gii.Name == gij.Name {
						return gii.Prov < gij.Prov
					}
					return gii.Name < gij.Name
				}
			})
			for _, gi := range gamelist {
				fmt.Println(FormatGameInfo(gi))
			}
		}
		if is, _ := listflags.GetBool("stat"); is {
			var provlist = make([]string, 0, len(prov))
			for p, n := range prov {
				provlist = append(provlist, fmt.Sprintf("%s: %d games", p, n))
			}

			fmt.Println()
			fmt.Printf("total: %d games, %d algorithms, %d providers\n", len(gamelist), len(alg), len(prov))
			sort.Strings(provlist)
			for _, s := range provlist {
				fmt.Println(s)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listflags = listCmd.Flags()
	listflags.BoolP("name", "n", true, "list of provided games names")
	listflags.BoolP("stat", "t", true, "summary statistics of provided games")

	listflags.BoolVarP(&fSort, "sort", "s", false, "sort by provider, else sort by names")
	listflags.BoolVarP(&fProp, "prop", "p", false, "print properties for each game")
	listflags.Float64Var(&fMrtp, "mrtp", 0, "RTP (Return to Player) of reels closest to given master RTP")
	listflags.Float64Var(&fDiff, "diff", 0, "difference between master RTP and closest to it real reels RTP")
	listflags.BoolVarP(&fRTP, "rtp", "r", false, "RTP (Return to Player) percents list of available reels for each game")

	listflags.StringSliceVarP(&inclist, "include", "i", []string{"all"}, "filter(s) to include games, filters can be as follows:\n"+
		"slot - all slot games\n"+
		"keno - all keno games\n"+
		"agt - games of 'AGT' provider\n"+
		"aristocrat - games of 'Aristocrat' provider\n"+
		"betsoft - games of 'BetSoft' provider\n"+
		"igt - games of 'IGT' provider\n"+
		"megajack - games of 'Megajack' provider\n"+
		"netent - games of 'NetExt' provider\n"+
		"novomatic - games of 'Novomatic' provider\n"+
		"playngo - games of 'Play'n GO' provider\n"+
		"playtech - games of 'Playtech' provider\n"+
		"slotopol - games of this 'Slotopol' provider\n"+
		"3x, 4x, 5x, ... - games with 3, 4, 5, ... reels\n"+
		"3x3, 4x4, 5x3, ... - games with 3x3, 4x4, 5x3, ... screen dimension\n"+
		"lines - games with wins counted by lines\n"+
		"ln=10 - games with 10 bet lines (or some other pointed)\n"+
		"ln<10 - games with less than 10 bet lines (or some other pointed)\n"+
		"ln>10 - games with greater than 10 bet lines (or some other pointed)\n"+
		"ways - games with wins counted by multiways, i.e. with 243, 1024 ways\n"+
		"wn=243 - games with 243 ways (or some other pointed)\n"+
		"wn<243 - games with less than 243 ways (or some other pointed)\n"+
		"wn>243 - games with greater than 243 ways (or some other pointed)\n"+
		"casc - slots with cascade falls\n"+
		"jack - games with jackpots\n"+
		"fg - games with any free games\n"+
		"bon - games with bonus games\n"+
		"y=15, y=2015 - games released in 2015 year (or some other pointed year)\n"+
		"y<15, y<2015 - games released before 2015 year (or some other pointed year)\n"+
		"y>15, y>2015 - games released after 2015 year (or some other pointed year)\n"+
		"all - all games")
	listflags.StringSliceVarP(&exclist, "exclude", "e", nil, "filter(s) to exclude games, filters are same as for include option")

	listflags.SortFlags = false
}
