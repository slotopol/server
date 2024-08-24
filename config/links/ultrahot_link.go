//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/ultrahot"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"ultrahot", "Ultra Hot"},
		},
		Provider: "Novomatic",
		ScrnX:    3,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(ultrahot.ReelsMap))
	for rtp := range ultrahot.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				ultrahot.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return ultrahot.NewGame()
		}
	}
}
