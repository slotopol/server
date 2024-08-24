//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/beetlemania"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"beetlemania", "Beetle Mania"},
			{"beetlemaniadeluxe", "Beetle Mania Deluxe"},
			{"hottarget", "Hot Target"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(beetlemania.ReelsMap))
	for rtp := range beetlemania.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" || rn == "bonu" {
					beetlemania.CalcStatBon(ctx, rn)
				} else {
					beetlemania.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return beetlemania.NewGame()
		}
	}
}
