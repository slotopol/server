//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/columbus"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("columbus", false, "'Columbus' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
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
		cfg.GameAliases[alias] = "columbus"
	}

	cfg.GameFactory["columbus"] = func(rd string) any {
		if _, ok := columbus.ReelsMap[rd]; ok {
			return columbus.NewGame(rd)
		}
		return nil
	}
}
