//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/plentyontwenty"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"plentyontwenty", "Plenty on Twenty"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(plentyontwenty.ReelsMap))
	for rtp := range plentyontwenty.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				plentyontwenty.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return plentyontwenty.NewGame(rtp)
		}
	}
}
