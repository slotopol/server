//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/beetlemania"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("beetlemania", false, "'Beetle Mania' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("beetlemania"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" || rn == "bonu" {
				beetlemania.CalcStatBon(ctx, rn)
			} else {
				beetlemania.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"beetlemania",
		"beetlemaniadeluxe",
		"hottarget",
	} {
		cfg.GameAliases[alias] = "beetlemania"
	}

	cfg.GameFactory["beetlemania"] = func(rd string) any {
		if _, ok := beetlemania.ReelsMap[rd]; ok {
			return beetlemania.NewGame(rd)
		}
		return nil
	}
}
