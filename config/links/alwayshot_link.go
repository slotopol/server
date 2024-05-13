//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/alwayshot"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("alwayshot", false, "'Always Hot' Novomatic 3x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("alwayshot"); is {
			var rn, _ = flags.GetString("reels")
			alwayshot.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"alwayshot",
	} {
		GameAliases[alias] = "alwayshot"
	}

	GameFactory["alwayshot"] = func(rd string) any {
		if _, ok := alwayshot.ReelsMap[rd]; ok {
			return alwayshot.NewGame(rd)
		}
		return nil
	}
}
