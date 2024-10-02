//go:build !prod || full || netent

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/simsalabim"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "simsalabim", Name: "Simsalabim"},
		},
		Provider: "NetEnt",
		SX:       5,
		SY:       3,
		LN:       25,
		FG:       FGretrig,
		BN:       1,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					slot.CalcStatBon(ctx)
				} else {
					slot.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
