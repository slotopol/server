//go:build !prod || full || betsoft

package links

import (
	"context"

	"github.com/slotopol/server/game/atthemovies"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("atthemovies", false, "'At the Movies' BetSoft 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("atthemovies"); is {
			var rn, _ = flags.GetString("reels")
			atthemovies.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"atthemovies",
	} {
		GameAliases[alias] = "atthemovies"
	}

	GameFactory["atthemovies"] = func(rd string) any {
		if _, ok := atthemovies.ReelsMap[rd]; ok {
			return atthemovies.NewGame(rd)
		}
		return nil
	}
}
