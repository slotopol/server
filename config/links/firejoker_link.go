//go:build !prod || full || playngo

package links

import (
	"context"

	"github.com/slotopol/server/game/firejoker"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"firejoker", "Fire Joker"},
		},
		Provider: "Play'n GO",
		ScrnX:    5,
		ScrnY:    3,
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				firejoker.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := firejoker.ReelsMap[rd]; ok {
				return firejoker.NewGame(rd)
			}
			return nil
		}
	}
}
