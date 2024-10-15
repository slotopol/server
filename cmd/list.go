package cmd

import (
	"fmt"
	"sort"
	"strings"

	cfg "github.com/slotopol/server/config"
	game "github.com/slotopol/server/game"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var listflags *pflag.FlagSet

var (
	fAll, fProp, fRTP bool
)

const listShort = "List of available games released on server"
const listLong = ``
const listExmp = `Get the list of all available games:
  %[1]s list --all
Get the list of available 'NetExt' and 'BetSoft' games only:
  %[1]s list --netent --betsoft
Get the list of available 'Play'n GO' games with RTP list for each:
  %[1]s list --playngo --rtp`

func Include(gi *game.GameInfo) bool {
	if fAll {
		return true
	}
	var is bool
	if is, _ = listflags.GetBool(util.ToID(gi.Provider)); is {
		return true
	}
	if is, _ = listflags.GetBool("3reels"); is && gi.SX == 3 {
		return true
	}
	if is, _ = listflags.GetBool("5reels"); is && gi.SX == 5 {
		return true
	}
	if is, _ = listflags.GetBool("3x3"); is && gi.SX == 3 && gi.SY == 3 {
		return true
	}
	if is, _ = listflags.GetBool("5x3"); is && gi.SX == 5 && gi.SY == 3 {
		return true
	}
	if is, _ = listflags.GetBool("5x4"); is && gi.SX == 5 && gi.SY == 4 {
		return true
	}
	if is, _ = listflags.GetBool("fewlines"); is && gi.LN < 20 {
		return true
	}
	if is, _ = listflags.GetBool("multilines"); is && gi.LN >= 20 {
		return true
	}
	if is, _ = listflags.GetBool("megaway"); is && gi.LN > 100 {
		return true
	}
	if is, _ = listflags.GetBool("fg"); is && gi.GP&(game.GPfghas+game.GPretrig) > 0 {
		return true
	}
	if is, _ = listflags.GetBool("bonus"); is && gi.BN > 0 {
		return true
	}
	return false
}

func FormatGameInfo(gi *game.GameInfo, ai int) string {
	var buf = make([]string, 0, 10)
	if gi.SN > 0 {
		buf = append(buf, fmt.Sprintf("'%s' %s %dx%d videoslot", gi.Aliases[ai].Name, gi.Provider, gi.SX, gi.SY))
	} else {
		buf = append(buf, fmt.Sprintf("'%s' %s %d spots lottery", gi.Aliases[ai].Name, gi.Provider, gi.SX))
	}
	if fProp {
		if gi.SN > 0 {
			buf = append(buf, fmt.Sprintf("%d symbols", gi.SN))
		}
		if gi.LN > 100 {
			buf = append(buf, fmt.Sprintf("%d ways", gi.LN))
		} else if gi.LN > 0 {
			var s string
			if gi.GP&game.GPsel == 0 {
				s = "constant "
			}
			buf = append(buf, fmt.Sprintf("%s%d lines", s, gi.LN))
		}
		if gi.GP&game.GPjack > 0 {
			buf = append(buf, "has jackpot")
		}
		if gi.GP&game.GPscat > 0 {
			buf = append(buf, "has scatters")
		}
		if gi.GP&(game.GPfghas+game.GPretrig) > 0 {
			var s1, s2, s3 string
			if gi.GP&game.GPretrig > 0 {
				s1 = "retriggerable "
			}
			if gi.GP&game.GPfgmult > 0 {
				s2 = " with multiplier"
			}
			if gi.GP&game.GPfgreel > 0 {
				s3 = " on bonus reels"
			}
			buf = append(buf, s1+"free games"+s2+s3)
		}
		if gi.GP&game.GPwild > 0 {
			buf = append(buf, "has wilds")
		}
		if gi.GP&game.GPrwild > 0 {
			buf = append(buf, "has reel wilds")
		}
		if gi.GP&game.GPbwild > 0 {
			buf = append(buf, "has big wilds")
		}
	}
	if fRTP && len(gi.RTP) > 0 {
		var rtpbuf = make([]string, len(gi.RTP))
		for i, rtp := range gi.RTP {
			rtpbuf[i] = fmt.Sprintf("%.2f", rtp)
		}
		buf = append(buf, "RTP: "+strings.Join(rtpbuf, ", "))
	}
	return strings.Join(buf, ", ")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Long:    listLong,
	Example: fmt.Sprintf(listExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var num, alg int
		var prov = map[string]int{}
		for _, gi := range game.GameList {
			if Include(gi) {
				num += len(gi.Aliases)
			}
		}

		var gamelist = make([]string, num)
		var i int
		for _, gi := range game.GameList {
			if Include(gi) {
				prov[gi.Provider] += len(gi.Aliases)
				alg++
				for ai := range gi.Aliases {
					gamelist[i] = FormatGameInfo(gi, ai)
					i++
				}
			}
		}
		var provlist = make([]string, len(prov))
		i = 0
		for p, n := range prov {
			provlist[i] = fmt.Sprintf("%s: %d games", p, n)
			i++
		}

		if is, _ := listflags.GetBool("name"); is {
			fmt.Println()
			sort.Strings(gamelist)
			for _, s := range gamelist {
				fmt.Println(s)
			}
		}
		if is, _ := listflags.GetBool("stat"); is {
			fmt.Println()
			fmt.Printf("total: %d games, %d algorithms, %d providers\n", num, alg, len(prov))
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
	listflags.BoolP("stat", "s", true, "summary statistics of provided games")

	listflags.BoolVar(&fAll, "all", false, "include all provided games, overrides any other filters")
	listflags.BoolVar(&fProp, "prop", false, "print properties for each game")
	listflags.BoolVar(&fRTP, "rtp", false, "RTP (Return to Player) percents list of available reels for each game")

	listflags.Bool("agt", false, "include games of 'AGT' provider")
	listflags.Bool("aristocrat", false, "include games of 'Aristocrat' provider")
	listflags.Bool("betsoft", false, "include games of 'BetSoft' provider")
	listflags.Bool("megajack", false, "include games of 'Megajack' provider")
	listflags.Bool("netent", false, "include games of 'NetExt' provider")
	listflags.Bool("novomatic", false, "include games of 'Novomatic' provider")
	listflags.Bool("playngo", false, "include games of 'Play'n GO' provider")
	listflags.Bool("playtech", false, "include games of 'Playtech' provider")
	listflags.Bool("slotopol", false, "include games of this 'Slotopol' provider")

	listflags.Bool("3reels", false, "include games with 3 reels")
	listflags.Bool("5reels", false, "include games with 5 reels")
	listflags.Bool("3x3", false, "include games with 3x3 screen")
	listflags.Bool("5x3", false, "include games with 5x3 screen")
	listflags.Bool("5x4", false, "include games with 5x4 screen")
	listflags.Bool("fewlines", false, "include games with few lines, i.e. with less than 20")
	listflags.Bool("multilines", false, "include games with few lines, i.e. with not less than 20")
	listflags.Bool("megaway", false, "include games with multiways, i.e. with 243, 1024 ways")
	listflags.Bool("fg", false, "include games with any free games")
	listflags.Bool("bonus", false, "include games with bonus games")

	listCmd.MarkFlagsOneRequired("all",
		"agt", "aristocrat", "betsoft", "megajack", "netent", "novomatic", "playngo", "playtech", "slotopol",
		"3reels", "5reels", "3x3", "5x3", "5x4", "fewlines", "multilines", "megaway", "fg", "bonus")
}
