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
		RtpList:  MakeRtpList(copsnrobbers.ReelsMap),
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
		GameFactory[ga.ID] = func() any {
			return copsnrobbers.NewGame()
		}
	}
}
