//go:build !prod || full || megajack

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/slotopol"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.BoolP("slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("slotopol"); is {
			var rn, _ = flags.GetString("reels")
			slotopol.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"slotopol",
	} {
		cfg.GameAliases[alias] = "slotopol"
	}

	cfg.GameFactory["slotopol"] = func(rd string) any {
		if _, ok := slotopol.ReelsMap[rd]; ok {
			return slotopol.NewGame(rd)
		}
		return nil
	}
}
