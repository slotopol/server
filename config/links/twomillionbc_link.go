//go:build !prod || full || betsoft

package links

import (
	"context"

	"github.com/slotopol/server/game/twomillionbc"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"2millionbc", "2 Million B.C."},
		},
		Provider: "BetSoft",
		ScrnX:    5,
		ScrnY:    3,
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					twomillionbc.CalcStatBon(ctx)
				} else {
					twomillionbc.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := twomillionbc.ReelsMap[rd]; ok {
				return twomillionbc.NewGame(rd)
			}
			return nil
		}
	}
}
