//go:build !prod || full || netent

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/groovysixties"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "groovysixties", Name: "Groovy Sixties"},
			{ID: "funkyseventies", Name: "Funky Seventies"}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
			{ID: "supereighties", Name: "Super Eighties"},   // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
		},
		Provider: "NetEnt",
		SX:       5,
		SY:       4,
		GP:       GPsel | GPretrig | GPfgmult | GPscat | GPwild,
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
