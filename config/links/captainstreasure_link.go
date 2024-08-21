//go:build !prod || full || playtech

package links

import (
	"context"

	"github.com/slotopol/server/game/captainstreasure"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"captainstreasure", "Captain's Treasure"},
		},
		Provider: "Playtech",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(captainstreasure.ReelsMap))
	for rtp := range captainstreasure.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				captainstreasure.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return captainstreasure.NewGame(rtp)
		}
	}
}
