//go:build !prod || full || megajack

package links

import (
	"context"

	"github.com/slotopol/server/game/slotopoldeluxe"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"slotopoldeluxe", "Slotopol Deluxe"},
		},
		Provider: "Megajack",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(slotopoldeluxe.ReelsMap))
	for rtp := range slotopoldeluxe.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				slotopoldeluxe.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slotopoldeluxe.NewGame()
		}
	}
}
