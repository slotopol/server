//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/katana"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"katana", "Katana"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(katana.ReelsMap))
	for rtp := range katana.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					katana.CalcStatBon(ctx)
				} else {
					katana.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return katana.NewGame(rtp)
		}
	}
}
