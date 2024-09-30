package cmd

import (
	"fmt"
	"sort"
	"strings"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/config/links"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var listflags *pflag.FlagSet

var (
	fAll, fRTP bool
)

const listShort = "List of available games released on server"
const listLong = ``
const listExmp = `Get the list of all available games:
  %[1]s list --all
Get the list of available 'NetExt' and 'BetSoft' games only:
  %[1]s list --netent --betsoft
Get the list of available 'Play'n GO' games with RTP list for each:
  %[1]s list --playngo --rtp`

func Include(gi *links.GameInfo) bool {
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
	if is, _ = listflags.GetBool("fg"); is && gi.FG >= links.FGhas {
		return true
	}
	if is, _ = listflags.GetBool("bonus"); is && gi.BN > 0 {
		return true
	}
	return false
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
		for _, gi := range links.GameList {
			if Include(&gi) {
				num += len(gi.Aliases)
			}
		}

		var gamelist = make([]string, num)
		var i int
		for _, gi := range links.GameList {
			if Include(&gi) {
				prov[gi.Provider] += len(gi.Aliases)
				alg++
				for _, ga := range gi.Aliases {
					var rtpinfo string
					if fRTP && len(gi.RTP) > 0 {
						var rtpstr = make([]string, len(gi.RTP))
						for i, rtp := range gi.RTP {
							rtpstr[i] = fmt.Sprintf("%.2f", rtp)
						}
						rtpinfo = ", RTP: " + strings.Join(rtpstr, ", ")
					}
					if gi.LN > 100 {
						gamelist[i] = fmt.Sprintf("'%s' %s %dx%d videoslot, %d ways%s", ga.Name, gi.Provider, gi.SX, gi.SY, gi.LN, rtpinfo)
					} else if gi.SY > 0 {
						gamelist[i] = fmt.Sprintf("'%s' %s %dx%d videoslot, %d lines%s", ga.Name, gi.Provider, gi.SX, gi.SY, gi.LN, rtpinfo)
					} else {
						gamelist[i] = fmt.Sprintf("'%s' %s %d spots lottery%s", ga.Name, gi.Provider, gi.SX, rtpinfo)
					}
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
	listflags.BoolVar(&fRTP, "rtp", false, "RTP (Return to Player) percents list of available reels for each game")

	listflags.Bool("aristocrat", false, "include games of 'Aristocrat' provider")
	listflags.Bool("megajack", false, "include games of 'Megajack' provider")
	listflags.Bool("novomatic", false, "include games of 'Novomatic' provider")
	listflags.Bool("netent", false, "include games of 'NetExt' provider")
	listflags.Bool("betsoft", false, "include games of 'BetSoft' provider")
	listflags.Bool("playtech", false, "include games of 'Playtech' provider")
	listflags.Bool("playngo", false, "include games of 'Play'n GO' provider")
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
		"aristocrat", "megajack", "novomatic", "netent", "betsoft", "playtech", "playngo", "slotopol",
		"3reels", "5reels", "3x3", "5x3", "5x4", "fewlines", "multilines", "megaway", "fg", "bonus")
}
