//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/arabiannights"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("arabiannights", false, "'Arabian Nights' NetEnt 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("arabiannights"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				arabiannights.CalcStatBon(ctx)
			} else {
				arabiannights.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"arabiannights",
	} {
		GameAliases[alias] = "arabiannights"
	}

	GameFactory["arabiannights"] = func(rd string) any {
		if _, ok := arabiannights.ReelsMap[rd]; ok {
			return arabiannights.NewGame(rd)
		}
		return nil
	}
}
