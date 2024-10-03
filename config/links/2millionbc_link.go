//go:build !prod || full || betsoft

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/2millionbc"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "2millionbc", Name: "2 Million B.C."},
		},
		Provider: "BetSoft",
		SX:       5,
		SY:       3,
		GP:       GPsel | GPretrig | GPfgreel | GPscat,
		SN:       len(slot.LinePay),
		LN:       30,
		BN:       2,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					slot.CalcStatBon(ctx)
				} else {
					slot.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
