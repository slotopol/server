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
		RtpList: []string{
			"93", "94", "95", "97", "100",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				atthemovies.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := atthemovies.ReelsMap[rd]; ok {
				return atthemovies.NewGame(rd)
			}
			return nil
		}
	}
}
