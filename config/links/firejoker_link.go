//go:build !prod || full || playngo

package links

import (
	"context"

	"github.com/slotopol/server/game/firejoker"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("firejoker", false, "'Fire Joker' Play'n GO 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("firejoker"); is {
			var rn, _ = flags.GetString("reels")
			firejoker.CalcStatReg(ctx, rn)
		}
	})

	for _, alias := range []string{
		"firejoker",
	} {
		GameAliases[alias] = "firejoker"
	}

	GameFactory["firejoker"] = func(rd string) any {
		if _, ok := firejoker.ReelsMap[rd]; ok {
			return firejoker.NewGame(rd)
		}
		return nil
	}
}
