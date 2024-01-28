//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/katana"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("katana", false, "'Katana' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("katana"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				katana.CalcStatBon(ctx)
			} else {
				katana.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"katana",
	} {
		cfg.GameAliases[alias] = "katana"
	}

	cfg.GameFactory["katana"] = func(rd string) any {
		if _, ok := katana.ReelsMap[rd]; ok {
			return katana.NewGame(rd)
		}
		return nil
	}
}
