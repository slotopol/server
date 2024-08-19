//go:build !prod || full || novomatic

package links

import (
	"context"
	"strconv"

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
		RtpList: []string{
			"89", "93", "94", "95", "97", "101", "111",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				roaringforties.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			var rtp, _ = strconv.ParseFloat(rd, 64)
			return roaringforties.NewGame(rtp)
		}
	}
}
