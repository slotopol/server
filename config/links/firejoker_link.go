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
	gi.RtpList = make([]float64, 0, len(firejoker.ReelsMap))
	for rtp := range firejoker.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				firejoker.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return firejoker.NewGame(rtp)
		}
	}
}
