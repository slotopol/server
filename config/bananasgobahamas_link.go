//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/bananasgobahamas"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("bananasgobahamas", false, "'Bananas Go Bahamas' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("bananasgobahamas"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				bananasgobahamas.CalcStatBon(ctx)
			} else {
				bananasgobahamas.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"bananasgobahamas",
	} {
		GameAliases[alias] = "bananasgobahamas"
	}

	GameFactory["bananasgobahamas"] = func(name string) any {
		if _, ok := bananasgobahamas.ReelsMap[name]; ok {
			return bananasgobahamas.NewGame(name)
		}
		return nil
	}
}
