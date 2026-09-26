package cmd

import (
	"fmt"
	"sort"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"

	"github.com/spf13/cobra"
)

const dataShort = "List of data routes for available games"
const dataLong = "List of data routes for available games. " +
	"Each route points to associated data at YAML-file. " +
	"In most common cases this is the reels sets data for provided games. " +
	"Data could be partially loaded by skip of embedded data load " +
	"(flag '--noembed') and by load data from external YAML-files (flags '-f')."
const dataExmp = `Get the list of routes with the loaded data:
  %[1]s data
Skip embedded data loading and get the list of loaded data from file 'reeldev.yaml':
  %[1]s --noembed -f=reeldev.yaml data
Show loaded data in the list of all available routes:
  %[1]s -vv --noembed -f=reeldev.yaml data -rl`

// dataCmd represents the `data` command
var dataCmd = &cobra.Command{
	Use:     "data",
	Short:   dataShort,
	Long:    dataLong,
	Example: fmt.Sprintf(dataExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var exitctx = Startup()
		var pf = cmd.Flags()

		// Load yaml-files
		var noembed bool
		if noembed, err = pf.GetBool("noembed"); err != nil {
			cfg.Fatalf(err.Error())
			return
		}
		if !noembed {
			if _, err = LoadInternalYaml(exitctx); err != nil {
				cfg.Fatalf("can not load internal yaml files: %s", err.Error())
				return
			}
		}
		if err = LoadExternalYaml(exitctx); err != nil {
			cfg.Fatalf("can not load external yaml files: %s", err.Error())
			return
		}
		UpdateAlgList()

		type item struct {
			has   bool
			route string
		}

		var routes bool
		if routes, err = pf.GetBool("routes"); err != nil {
			cfg.Fatalf(err.Error())
			return
		}
		var loaded bool
		if loaded, err = pf.GetBool("loaded"); err != nil {
			cfg.Fatalf(err.Error())
			return
		}

		var list = make([]item, 0, len(game.DataRouter))
		for id := range game.DataRouter {
			var _, ok = game.DataLoaded[id]
			list = append(list, item{has: ok, route: id})
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].route < list[j].route
		})
		for _, r := range list {
			if routes || r.has {
				if loaded {
					var c byte
					if r.has {
						c = '+'
					} else {
						c = '-'
					}
					fmt.Printf("[%c] %s\n", c, r.route)
				} else {
					fmt.Printf("%s\n", r.route)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dataCmd)

	var f = dataCmd.Flags()
	f.Bool("noembed", false, "do not load embedded yaml files, useful for development")
	f.BoolP("routes", "r", false, "print the list of all registered routes or otherwise the list of routes with loaded data only")
	f.BoolP("loaded", "l", false, "print the marker that points to data loaded status")

	f.SortFlags = false
}
