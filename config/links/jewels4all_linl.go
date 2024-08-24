//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/jewels4all"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"jewels4all", "Jewels 4 All"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(jewels4all.ChanceMap))
	for rtp := range jewels4all.ChanceMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				jewels4all.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return jewels4all.NewGame()
		}
	}
}
