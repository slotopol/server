//go:build !prod || full || betsoft

package links

import (
	"context"

	"github.com/slotopol/server/game/twomillionbc"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("2millionbc", false, "'2 Million B.C.' BetSoft 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("2millionbc"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				twomillionbc.CalcStatBon(ctx)
			} else {
				twomillionbc.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"2millionbc",
	} {
		GameAliases[alias] = "2millionbc"
	}

	GameFactory["2millionbc"] = func(rd string) any {
		if _, ok := twomillionbc.ReelsMap[rd]; ok {
			return twomillionbc.NewGame(rd)
		}
		return nil
	}
}
