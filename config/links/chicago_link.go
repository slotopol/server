//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/chicago"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"chicago", "Chicago"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(chicago.ReelsMap))
	for rtp := range chicago.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				chicago.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return chicago.NewGame()
		}
	}
}
