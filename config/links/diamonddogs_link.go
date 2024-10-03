//go:build !prod || full || netent

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/diamonddogs"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "diamonddogs", Name: "Diamond Dogs"},
			{ID: "voodoovibes", Name: "Voodoo Vibes"},
		},
		Provider: "NetEnt",
		SX:       5,
		SY:       3,
		GP:       GPsel | GPretrig | GPfgmult | GPfgreel | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       25,
		BN:       1,
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
