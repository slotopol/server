//go:build !prod || full || playtech

package links

import (
	"context"

	"github.com/slotopol/server/game/panthermoon"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"panthermoon", "Panther Moon"},
			{"safariheat", "Safari Heat"},
		},
		Provider: "Playtech",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(panthermoon.ReelsMap))
	for rtp := range panthermoon.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					panthermoon.CalcStatBon(ctx)
				} else {
					panthermoon.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return panthermoon.NewGame(rtp)
		}
	}
}
