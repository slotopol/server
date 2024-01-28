//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/jewels"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("jewels", false, "'Jewels' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("jewels"); is {
			var rn, _ = flags.GetString("reels")
			jewels.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"jewels",
	} {
		cfg.GameAliases[alias] = "jewels"
	}

	cfg.GameFactory["jewels"] = func(rd string) any {
		if _, ok := jewels.ReelsMap[rd]; ok {
			return jewels.NewGame(rd)
		}
		return nil
	}
}
