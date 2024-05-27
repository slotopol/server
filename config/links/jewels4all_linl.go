//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/jewels4all"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"jewels4all", "Jewels 4 All"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				jewels4all.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := jewels4all.ChanceMap[rd]; ok {
				return jewels4all.NewGame(rd)
			}
			return nil
		}
	}
}
