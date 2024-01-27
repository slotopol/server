//go:build !prod || full || megajack

package config

import (
	"context"

	"github.com/slotopol/server/game/slotopol"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.BoolP("slotopol", "s", false, "'Slotopol' Megajack 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("slotopol"); is {
			var rn, _ = flags.GetString("reels")
			slotopol.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"slotopol",
	} {
		GameAliases[alias] = "slotopol"
	}

	GameFactory["slotopol"] = func(rd string) any {
		if _, ok := slotopol.ReelsMap[rd]; ok {
			return slotopol.NewGame(rd)
		}
		return nil
	}
}
