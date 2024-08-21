//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/alwayshot"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"alwayshot", "Always Hot"},
		},
		Provider: "Novomatic",
		ScrnX:    3,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(alwayshot.ReelsMap))
	for rtp := range alwayshot.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				alwayshot.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return alwayshot.NewGame(rtp)
		}
	}
}
