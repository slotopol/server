//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/diamonddogs"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"diamonddogs", "Diamond Dogs"},
			{"voodoovibes", "Voodoo Vibes"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
		RtpList:  MakeRtpList(diamonddogs.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					diamonddogs.CalcStatBon(ctx)
				} else {
					diamonddogs.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return diamonddogs.NewGame()
		}
	}
}
