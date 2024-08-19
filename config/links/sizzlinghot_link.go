//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/sizzlinghot"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"sizzlinghot", "Sizzling Hot"},
			{"sizzlinghotdeluxe", "Sizzling Hot Deluxe"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"96",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				sizzlinghot.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			return sizzlinghot.NewGame(Atof(rd))
		}
	}
}
