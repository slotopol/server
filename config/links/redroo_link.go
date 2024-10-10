//go:build !prod || full || aristocrat

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/redroo"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "redroo", Name: "Redroo"},
		},
		Provider: "Aristocrat",
		SX:       5,
		SY:       4,
		GP:       GPretrig | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       1024,
		BN:       0,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var mrtp, _ = flags.GetFloat64("reels")
				slot.CalcStatReg(ctx, mrtp)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
