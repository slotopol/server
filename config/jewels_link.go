//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/jewels"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("jewels", false, "'Jewels' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("jewels"); is {
			var rn, _ = flags.GetString("reels")
			jewels.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"jewels",
	} {
		GameAliases[alias] = "jewels"
	}

	GameFactory["jewels"] = func(rd string) any {
		if _, ok := jewels.ReelsMap[rd]; ok {
			return jewels.NewGame(rd)
		}
		return nil
	}
}
