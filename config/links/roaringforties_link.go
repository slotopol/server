//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/roaringforties"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"roaringforties", "Roaring Forties"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    4,
		RtpList:  MakeRtpList(roaringforties.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				roaringforties.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return roaringforties.NewGame()
		}
	}
}
