//go:build full || megajack

package config

import (
	"context"

	"github.com/slotopol/server/game/slotopoldeluxe"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("slotopoldeluxe", false, "'Slotopol Deluxe' Megajack 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("slotopoldeluxe"); is {
			var rn, _ = flags.GetString("reels")
			slotopoldeluxe.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"slotopoldeluxe",
	} {
		GameAliases[alias] = "slotopoldeluxe"
	}

	GameFactory["slotopoldeluxe"] = func(name string) any {
		if reels, ok := slotopoldeluxe.ReelsMap[name]; ok {
			return slotopoldeluxe.NewGame(reels)
		}
		return nil
	}
}
