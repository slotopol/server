//go:build !prod || full || novomatic

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/sizzlinghot"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "sizzlinghot", Name: "Sizzling Hot"},
			{ID: "sizzlinghotdeluxe", Name: "Sizzling Hot Deluxe"},
		},
		Provider: "Novomatic",
		SX:       5,
		SY:       3,
		GP:       GPfgno,
		SN:       len(slot.LinePay),
		LN:       5,
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
