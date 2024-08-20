//go:build !prod || full || playngo

package links

import (
	"context"

	"github.com/slotopol/server/game/copsnrobbers"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"copsnrobbers", "Cops'n'Robbers"},
		},
		Provider: "Play'n GO",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"86", "88", "90", "92", "93", "94", "95", "96", "97", "98", "99", "112", "bon",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					copsnrobbers.CalcStatBon(ctx)
				} else {
					copsnrobbers.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return copsnrobbers.NewGame(rtp)
		}
	}
}
