//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/plentyontwenty"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("plentyontwenty", false, "'Plenty on Twenty' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("plentyontwenty"); is {
			var rn, _ = flags.GetString("reels")
			plentyontwenty.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"plentyontwenty",
	} {
		GameAliases[alias] = "plentyontwenty"
	}

	GameFactory["plentyontwenty"] = func(name string) any {
		if _, ok := plentyontwenty.ReelsMap[name]; ok {
			return plentyontwenty.NewGame(name)
		}
		return nil
	}
}
