//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/trolls"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"trolls", "Trolls"},
			{"excalibur", "Excalibur"},
			{"pandorasbox", "Pandora's Box"},
			{"wildwitches", "Wild Witches"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
	}
	gi.RtpList = make([]float64, 0, len(trolls.ReelsMap))
	for rtp := range trolls.ReelsMap {
		gi.RtpList = append(gi.RtpList, rtp)
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				trolls.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return trolls.NewGame(rtp)
		}
	}
}
