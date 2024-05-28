//go:build !prod || full || novomatic

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
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"86", "88", "90", "92", "94", "95", "96", "97", "141", "bon",
		},
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
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := panthermoon.ReelsMap[rd]; ok {
				return panthermoon.NewGame(rd)
			}
			return nil
		}
	}
}
