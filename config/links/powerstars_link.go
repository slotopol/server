//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/powerstars"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("powerstars", false, "'Power Stars' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("powerstars"); is {
			var rn, _ = flags.GetString("reels")
			powerstars.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"powerstars",
	} {
		GameAliases[alias] = "powerstars"
	}

	GameFactory["powerstars"] = func(rd string) any {
		if _, ok := powerstars.ChanceMap[rd]; ok {
			return powerstars.NewGame(rd)
		}
		return nil
	}
}
