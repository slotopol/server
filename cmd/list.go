package cmd

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"

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
Get the list of 'AGT' games with screen 3x3 only:
  %[1]s list -i agt+3x3
Get the list of 'AGT' games with big symbols and free games:
  %[1]s list -i agt+big+fg
Get the list of slots without scatters with more than 3 reels:
  %[1]s list -e scat -e keno -e 3x
  %[1]s list -i slot+~scat+~3x
Get the list of Megajack games with properties and RTP list for each:
  %[1]s list -i megajack --prop --rtp`

var listflags *pflag.FlagSet

var (
	fSort            bool
	fProp, fRTP      bool
	fMrtp, fDiff     float64
	inclist, exclist []string
)

func FormatGameInfo(gi *game.GameInfo) string {
	var b strings.Builder
	if gi.SN > 0 {
		if gi.GP&(game.GPcpay+game.GPcasc) == game.GPcpay+game.GPcasc {
			fmt.Fprintf(&b, "'%s' %s %dx%d cluster cascade videoslot", gi.Name, gi.Prov, gi.SX, gi.SY)
		} else if gi.GP&game.GPcpay > 0 {
			fmt.Fprintf(&b, "'%s' %s %dx%d cluster videoslot", gi.Name, gi.Prov, gi.SX, gi.SY)
		} else if gi.GP&game.GPcasc > 0 {
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
		if gi.LNum > 0 {
			if gi.GP&game.GPlsel == 0 {
				fmt.Fprintf(&b, ", constant %d lines", gi.LNum)
			} else {
				fmt.Fprintf(&b, ", %d lines", gi.LNum)
			}
		}
		if gi.WN > 0 {
			fmt.Fprintf(&b, ", %d ways", gi.WN)
		}
		if gi.BN > 0 {
			fmt.Fprintf(&b, ", %d bonus games", gi.BN)
		}
		if gi.GP&game.GPjack > 0 {
			b.WriteString(", has jackpot")
		}
		if gi.GP&game.GPfill > 0 {
			b.WriteString(", has multiplier on filled screen")
		}
		if gi.GP&game.GPcumul > 0 {
			b.WriteString(", has cumulative pays")
		}
		if gi.GP&game.GPbmode > 0 {
			b.WriteString(", has non-reels bonus mode")
		}
		if gi.GP&(game.GPfgany) > 0 {
			b.WriteString(", ")
			if gi.GP&game.GPfgseq > 0 {
				b.WriteString("retriggerable ")
			} else if gi.GP&game.GPfgtwic > 0 {
				b.WriteString("once retriggerable ")
			}
			b.WriteString("free games")
			if gi.GP&game.GPfgmult > 0 {
				b.WriteString(" with multiplier")
			}
			if gi.GP&game.GPfgreel > 0 {
				b.WriteString(" on bonus reels")
			}
		}
		if gi.GP&game.GPwsc > 0 {
			b.WriteString(", has wild/scatters")
		} else if gi.GP&game.GPscany > 0 {
			b.WriteString(", has scatters")
		}
		if gi.GP&game.GPwany > 0 {
			if gi.GP&game.GPwild > 0 {
				b.WriteString(", has wilds")
			}
			if gi.GP&game.GPrwild > 0 {
				b.WriteString(", has reel wilds")
			}
			if gi.GP&game.GPbwild > 0 {
				b.WriteString(", has big wilds")
			}
			if gi.GP&game.GPewild > 0 {
				b.WriteString(", has expanding wilds")
			}
			if gi.GP&game.GPwmult > 0 {
				b.WriteString(" with multiplier")
			}
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
		var err error
		var exitctx = context.Background()

		// Load yaml-files
		if fRTP || fMrtp > 0 || fDiff > 0 {
			LoadInternalYaml(exitctx)
			if err = LoadExternalYaml(exitctx); err != nil {
				log.Fatalf("can not load yaml files: %s", err.Error())
				return
			}
			UpdateAlgList()
			CheckAlgList()
		}

		var finclist, fexclist [][]game.Filter
		var f game.Filter
		var flist []game.Filter
		for _, inc := range inclist {
			if inc == "" {
				continue
			}
			var keys = strings.Split(inc, "+")
			flist = nil
			for _, key := range keys {
				if f = game.GetFilter(key); f == nil {
					fmt.Printf("filter with name '%s' does not recognized\n", key)
					continue
				}
				flist = append(flist, f)
			}
			finclist = append(finclist, flist)
		}
		for _, exc := range exclist {
			if exc == "" {
				continue
			}
			var keys = strings.Split(exc, "+")
			flist = nil
			for _, key := range keys {
				if f = game.GetFilter(key); f == nil {
					fmt.Printf("filter with name '%s' does not recognized\n", key)
					continue
				}
				flist = append(flist, f)
			}
			fexclist = append(fexclist, flist)
		}

		var alg = map[*game.AlgDescr]int{}
		var prov = map[string]int{}
		var gamelist = make([]*game.GameInfo, 0, 256)
		for aid, gi := range game.InfoMap {
			_ = aid
			if game.Passes(gi, finclist, fexclist) {
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

const filtdescr = `filter(s) to include games, filters could be conjuncted with '+' sign, for logical negation used '~' before filter, syntax can be as follows:
slot - all slot games
keno - all keno games
agt - games of 'AGT' provider
aristocrat - games of 'Aristocrat' provider
betsoft - games of 'BetSoft' provider
ct - games of 'CT Interactive' provider
igt - games of 'IGT' provider
megajack - games of 'Megajack' provider
netent - games of 'NetExt' provider
novomatic - games of 'Novomatic' provider
playngo - games of 'Play'n GO' provider
playtech - games of 'Playtech' provider
slotopol - games of this 'Slotopol' provider
3x, 4x, 5x, ... - games with 3, 4, 5, ... reels
3x3, 4x4, 5x3, ... - games with 3x3, 4x4, 5x3, ... screen dimension
lines - games with wins counted by lines
ln=10, lneq10 - games with 10 bet lines (or some other pointed)
ln<10, lnlt10 - games with less than 10 bet lines (or some other pointed)
ln>10, lngt10 - games with greater than 10 bet lines (or some other pointed)
ways - games with wins counted by multiways, i.e. with 243, 1024 ways
wn=243, wneq243 - games with 243 ways (or some other pointed)
wn<243, wnlt243 - games with less than 243 ways (or some other pointed)
wn>243, wngt243 - games with greater than 243 ways (or some other pointed)
bon - games with bonus games
lpay - pays left to right
rpay - pays left to right and right to left
apay - pays for combination at any position
cpay - pays by clusters
jack - slots with jackpots
fill - has multiplier on filled slot screen
bm - slots with non-reels bonus mode
casc - slots with cascade falls
cm - multipliers on cascade avalanche
fg - slots with any free games
fgo - slots with non-retriggered free games
fgt - slots with free games that can be retriggered only once
fgs - slots with free games that can be retriggered in a sequence
fgr - slots with separate reels on free games
fgm - slots with any multipliers on free games
sany - slots with any scatter symbols
scat - slots with regular scatters
wany - slots with any wild symbols
wild - slots with regular wild symbols
wsc - slots with wild/scatters symbols
rw - slots with reel wilds
bw - slots with big wilds (3x3)
ew - slots with expanding wilds
wt - symbols turns to wilds
wm - slots with multiplier on wilds
big - slots with big symbol (usually 3x3 in the center on free games)
y=15, y=2015, yeq15 - games released in 2015 year (or some other pointed year)
y<15, y<2015, ylt15 - games released before 2015 year (or some other pointed year)
y>15, y>2015, ygt15 - games released after 2015 year (or some other pointed year)
nodate - games with unknown release date
all - all games`

func init() {
	rootCmd.AddCommand(listCmd)

	listflags = listCmd.Flags()
	listflags.BoolP("name", "n", true, "list of provided games names")
	listflags.BoolP("stat", "s", true, "summary statistics of provided games")

	listflags.BoolVar(&fSort, "sort", false, "sort by provider, else sort by names")
	listflags.BoolVar(&fProp, "prop", false, "print properties for each game")
	listflags.Float64Var(&fMrtp, "mrtp", 0, "RTP (Return to Player) of reels closest to given master RTP")
	listflags.Float64Var(&fDiff, "diff", 0, "difference between master RTP and closest to it real reels RTP")
	listflags.BoolVarP(&fRTP, "rtp", "r", false, "RTP (Return to Player) percents list of available reels for each game")

	listflags.StringSliceVarP(&inclist, "include", "i", []string{"all"}, filtdescr)
	listflags.StringSliceVarP(&exclist, "exclude", "e", nil, "filter(s) to exclude games, filters are same as for include option")

	listflags.SortFlags = false
}
