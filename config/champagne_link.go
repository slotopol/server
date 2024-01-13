//go:build full || megajack

package config

import (
	"context"

	"github.com/slotopol/server/game/champagne"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("champagne", false, "'Champagne' Megajack 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("champagne"); is {
			var rn, _ = flags.GetString("reels")
			champagne.CalcStatReg(ctx, rn)
		}
	})

	for _, alias := range []string{
		"champagne",
	} {
		GameAliases[alias] = "champagne"
	}

	GameFactory["champagne"] = func(name string) any {
		if reels, ok := champagne.ReelsMap[name]; ok {
			return champagne.NewGame(reels)
		}
		return nil
	}
}
