//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/beetlemania"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("beetlemania", false, "'Beetle Mania' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
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
		GameAliases[alias] = "beetlemania"
	}

	GameFactory["beetlemania"] = func(name string) any {
		if _, ok := beetlemania.ReelsMap[name]; ok {
			return beetlemania.NewGame(name)
		}
		return nil
	}
}
