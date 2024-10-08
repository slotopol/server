//go:build !prod || full || playngo

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/fortuneteller"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "fortuneteller", Name: "Fortune Teller"},
		},
		Provider: "Play'n GO",
		SX:       5,
		SY:       3,
		GP:       GPsel | GPfghas | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       20,
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
