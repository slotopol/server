//go:build !prod || full || megajack

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

	GameFactory["slotopoldeluxe"] = func(rd string) any {
		if _, ok := slotopoldeluxe.ReelsMap[rd]; ok {
			return slotopoldeluxe.NewGame(rd)
		}
		return nil
	}
}
