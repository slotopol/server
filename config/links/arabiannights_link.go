//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/arabiannights"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"arabiannights", "Arabian Nights"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"85", "88", "90", "91", "92", "93", "95", "96", "97", "99", "107", "bon",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					arabiannights.CalcStatBon(ctx)
				} else {
					arabiannights.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return arabiannights.NewGame(rtp)
		}
	}
}
