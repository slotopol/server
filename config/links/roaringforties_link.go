//go:build !prod || full || novomatic

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/roaringforties"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "roaringforties", Name: "Roaring Forties"},
		},
		Provider: "Novomatic",
		SX:       5,
		SY:       4,
		GP:       GPsel | GPfgno | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       40,
		BN:       0,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var mrtp, _ = flags.GetFloat64("reels")
				slot.CalcStat(ctx, mrtp)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
