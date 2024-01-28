//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/justjewels"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("justjewels", false, "'Just Jewels' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("justjewels"); is {
			var rn, _ = flags.GetString("reels")
			justjewels.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"justjewels",
	} {
		cfg.GameAliases[alias] = "justjewels"
	}

	cfg.GameFactory["justjewels"] = func(rd string) any {
		if _, ok := justjewels.ReelsMap[rd]; ok {
			return justjewels.NewGame(rd)
		}
		return nil
	}
}
