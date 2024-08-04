package cmd

import (
	"fmt"
	"sort"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/config/links"
	"github.com/slotopol/server/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var listflags *pflag.FlagSet

var (
	fAllPrv bool
)

const listShort = "List of games"
const listLong = ``
const listExmp = `Get the list of all available games:
  %s list --all`

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
			if prv, _ := listflags.GetBool(util.ToID(gi.Provider)); prv || fAllPrv {
				num += len(gi.Aliases)
			}
		}

		var list = make([]string, num)
		var i int
		for _, gi := range links.GameList {
			if prv, _ := listflags.GetBool(util.ToID(gi.Provider)); prv || fAllPrv {
				prov[gi.Provider] += len(gi.Aliases)
				alg++
				for _, ga := range gi.Aliases {
					list[i] = fmt.Sprintf("'%s' %s %dx%d videoslot", ga.Name, gi.Provider, gi.ScrnX, gi.ScrnY)
					i++
				}
			}
		}

		if is, _ := listflags.GetBool("name"); is {
			sort.Strings(list)
			fmt.Println()
			for _, s := range list {
				fmt.Println(s)
			}
		}
		if is, _ := listflags.GetBool("stat"); is {
			fmt.Println()
			fmt.Printf("total: %d games, %d algorithms, %d providers\n", num, alg, len(prov))
			for p, n := range prov {
				fmt.Printf("%s: %d games\n", p, n)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listflags = listCmd.Flags()
	listflags.BoolP("name", "n", true, "list of provided games names")
	listflags.BoolP("stat", "s", true, "summary statistics of provided games")
	listflags.BoolVar(&fAllPrv, "all", false, "list games of all available providers")
	listflags.Bool("megajack", false, "include games of 'Megajack' provider")
	listflags.Bool("novomatic", false, "include games of 'Novomatic' provider")
	listflags.Bool("netent", false, "include games of 'NetExt' provider")
	listflags.Bool("betsoft", false, "include games of 'BetSoft' provider")
	listflags.Bool("playtech", false, "include games of 'Playtech' provider")
	listflags.Bool("playngo", false, "include games of 'Play'n GO' provider")
	listCmd.MarkFlagsOneRequired("all", "megajack", "novomatic", "netent", "betsoft", "playtech", "playngo")
}
