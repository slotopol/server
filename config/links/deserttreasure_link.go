//go:build !prod || full || playtech

package links

import (
	"context"

	"github.com/slotopol/server/game/deserttreasure"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"deserttreasure", "Desert Treasure"},
		},
		Provider: "Playtech",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(deserttreasure.ReelsMap))
	for rtp := range deserttreasure.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					deserttreasure.CalcStatBon(ctx)
				} else {
					deserttreasure.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return deserttreasure.NewGame()
		}
	}
}
