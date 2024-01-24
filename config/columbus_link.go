//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/columbus"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("columbus", false, "'Columbus' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("columbus"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				columbus.CalcStatBon(ctx)
			} else {
				columbus.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"columbus",
		"columbusdeluxe",
		"marcopolo",
	} {
		GameAliases[alias] = "columbus"
	}

	GameFactory["columbus"] = func(name string) any {
		if _, ok := columbus.ReelsMap[name]; ok {
			return columbus.NewGame(name)
		}
		return nil
	}
}
