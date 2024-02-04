//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/chicago"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("chicago", false, "'Chicago' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("chicago"); is {
			var rn, _ = flags.GetString("reels")
			chicago.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"chicago",
	} {
		GameAliases[alias] = "chicago"
	}

	GameFactory["chicago"] = func(rd string) any {
		if _, ok := chicago.ReelsMap[rd]; ok {
			return chicago.NewGame(rd)
		}
		return nil
	}
}
