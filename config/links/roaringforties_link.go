//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/roaringforties"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("roaringforties", false, "'Roaring Forties' Novomatic 5x4 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("roaringforties"); is {
			var rn, _ = flags.GetString("reels")
			roaringforties.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"roaringforties",
	} {
		GameAliases[alias] = "roaringforties"
	}

	GameFactory["roaringforties"] = func(rd string) any {
		if _, ok := roaringforties.ReelsMap[rd]; ok {
			return roaringforties.NewGame(rd)
		}
		return nil
	}
}
