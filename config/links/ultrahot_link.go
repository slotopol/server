//go:build !prod || full || novomatic

package links

import (
	"context"
	"strconv"

	"github.com/slotopol/server/game/ultrahot"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"ultrahot", "Ultra Hot"},
		},
		Provider: "Novomatic",
		ScrnX:    3,
		ScrnY:    3,
		RtpList: []string{
			"88", "90", "92", "93", "96", "98", "111",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				ultrahot.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			var rtp, _ = strconv.ParseFloat(rd, 64)
			return ultrahot.NewGame(rtp)
		}
	}
}
