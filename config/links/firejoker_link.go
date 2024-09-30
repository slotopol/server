//go:build !prod || full || playngo

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/firejoker"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "firejoker", Name: "Fire Joker"},
		},
		Provider: "Play'n GO",
		SX:       5,
		SY:       3,
		LN:       5,
		FG:       FGhas,
		BN:       0,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				slot.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
