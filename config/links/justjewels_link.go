//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/justjewels"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("justjewels", false, "'Just Jewels' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("justjewels"); is {
			var rn, _ = flags.GetString("reels")
			justjewels.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"justjewels",
	} {
		GameAliases[alias] = "justjewels"
	}

	GameFactory["justjewels"] = func(rd string) any {
		if _, ok := justjewels.ReelsMap[rd]; ok {
			return justjewels.NewGame(rd)
		}
		return nil
	}
}
