//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/powerstars"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"powerstars", "Power Stars"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"86", "88", "90", "91", "92", "94", "95", "96", "98", "100", "112",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				powerstars.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := powerstars.ChanceMap[rd]; ok {
				return powerstars.NewGame(rd)
			}
			return nil
		}
	}
}
