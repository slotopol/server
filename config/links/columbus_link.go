//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/columbus"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"columbus", "Columbus"},
			{"columbusdeluxe", "Columbus Deluxe"},
			{"marcopolo", "Marco Polo"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(columbus.ReelsMap))
	for rtp := range columbus.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					columbus.CalcStatBon(ctx)
				} else {
					columbus.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return columbus.NewGame()
		}
	}
}
