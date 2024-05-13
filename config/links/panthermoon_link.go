//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/panthermoon"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("panthermoon", false, "'Panther Moon' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("panthermoon"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				panthermoon.CalcStatBon(ctx)
			} else {
				panthermoon.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"panthermoon",
		"safariheat",
	} {
		GameAliases[alias] = "panthermoon"
	}

	GameFactory["panthermoon"] = func(rd string) any {
		if _, ok := panthermoon.ReelsMap[rd]; ok {
			return panthermoon.NewGame(rd)
		}
		return nil
	}
}
