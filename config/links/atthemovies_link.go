//go:build !prod || full || betsoft

package links

import (
	"context"

	"github.com/slotopol/server/game/atthemovies"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"atthemovies", "At the Movies"},
		},
		Provider: "BetSoft",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(atthemovies.ReelsMap))
	for rtp := range atthemovies.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				atthemovies.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return atthemovies.NewGame()
		}
	}
}
