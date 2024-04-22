//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/jewels4all"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("jewels4all", false, "'Jewels 4 All' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("jewels4all"); is {
			var rn, _ = flags.GetString("reels")
			jewels4all.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"jewels4all",
	} {
		GameAliases[alias] = "jewels4all"
	}

	GameFactory["jewels4all"] = func(rd string) any {
		if _, ok := jewels4all.ChanceMap[rd]; ok {
			return jewels4all.NewGame(rd)
		}
		return nil
	}
}
