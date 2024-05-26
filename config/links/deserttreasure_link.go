//go:build !prod || full || playtech

package links

import (
	"context"

	"github.com/slotopol/server/game/deserttreasure"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("deserttreasure", false, "'Desert Treasure' Playtech 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("deserttreasure"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				deserttreasure.CalcStatBon(ctx)
			} else {
				deserttreasure.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"deserttreasure",
	} {
		GameAliases[alias] = "deserttreasure"
	}

	GameFactory["deserttreasure"] = func(rd string) any {
		if _, ok := deserttreasure.ReelsMap[rd]; ok {
			return deserttreasure.NewGame(rd)
		}
		return nil
	}
}
