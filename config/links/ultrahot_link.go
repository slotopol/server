//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/ultrahot"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("ultrahot", false, "'Ultra Hot' Novomatic 3x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("ultrahot"); is {
			var rn, _ = flags.GetString("reels")
			ultrahot.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"ultrahot",
		"ultrahotdeluxe",
	} {
		GameAliases[alias] = "ultrahot"
	}

	GameFactory["ultrahot"] = func(rd string) any {
		if _, ok := ultrahot.ReelsMap[rd]; ok {
			return ultrahot.NewGame(rd)
		}
		return nil
	}
}
