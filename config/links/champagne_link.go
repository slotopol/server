//go:build !prod || full || megajack

package links

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

	GameFactory["champagne"] = func(rd string) any {
		if _, ok := champagne.ReelsMap[rd]; ok {
			return champagne.NewGame(rd)
		}
		return nil
	}
}
