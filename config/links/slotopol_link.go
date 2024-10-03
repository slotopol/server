//go:build !prod || full || megajack

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/slotopol"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "slotopol", Name: "Slotopol"},
		},
		Provider: "Megajack",
		SX:       5,
		SY:       3,
		GP:       GPsel | GPfgno | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       21,
		BN:       2,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				slot.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
