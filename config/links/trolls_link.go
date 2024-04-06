//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/trolls"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("trolls", false, "'Trolls' NetEnt 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("trolls"); is {
			var rn, _ = flags.GetString("reels")
			trolls.CalcStatReg(ctx, rn)
		}
	})

	for _, alias := range []string{
		"trolls",
		"excalibur",
		"pandorasbox",
		"wildwitches",
	} {
		GameAliases[alias] = "trolls"
	}

	GameFactory["trolls"] = func(rd string) any {
		if _, ok := trolls.ReelsMap[rd]; ok {
			return trolls.NewGame(rd)
		}
		return nil
	}
}
