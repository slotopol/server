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
		RtpList: []string{
			"87", "89", "91", "93", "95", "97", "100", "bon",
		},
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
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := katana.ReelsMap[rd]; ok {
				return katana.NewGame(rd)
			}
			return nil
		}
	}
}
