//go:build !prod || full || novomatic

package config

import (
	"context"

	"github.com/slotopol/server/game/dolphinspearl"
	"github.com/spf13/pflag"
)

func init() {
	FlagsSetters = append(FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("dolphinspearl", false, "'Dolphins Pearl' Novomatic 5x3 slots")
	})
	ScatIters = append(ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("dolphinspearl"); is {
			var rn, _ = flags.GetString("reels")
			if rn == "bon" {
				dolphinspearl.CalcStatBon(ctx)
			} else {
				dolphinspearl.CalcStatReg(ctx, rn)
			}
		}
	})

	for _, alias := range []string{
		"dolphinspearl",
		"dolphinspearldeluxe",
		"attila",
		"bananasplash",
		"dynastyofming",
		"gryphonsgold",
		"jokerdolphin",
		"kingofcards",
		"luckyladyscharm",
		"luckyladyscharmdeluxe",
		"pharaonsgold2",
		"pharaonsgold3",
		"polarfox",
		"royaltreasures",
		"secretforest",
		"themoneygame",
		"unicornmagic",
	} {
		GameAliases[alias] = "dolphinspearl"
	}

	GameFactory["dolphinspearl"] = func(rd string) any {
		if _, ok := dolphinspearl.ReelsMap[rd]; ok {
			return dolphinspearl.NewGame(rd)
		}
		return nil
	}
}
