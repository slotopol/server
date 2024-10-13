//go:build !prod || full || netent

package groovysixties

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

func init() {
	var gi = game.GameInfo{
		Aliases: []game.GameAlias{
			{ID: "groovysixties", Name: "Groovy Sixties"},
			{ID: "funkyseventies", Name: "Funky Seventies"}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
			{ID: "supereighties", Name: "Super Eighties"},   // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
		},
		Provider: "NetEnt",
		SX:       5,
		SY:       4,
		GP: game.GPsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	}
	game.GameList = append(game.GameList, gi)

	for _, ga := range gi.Aliases {
		game.ScanIters = append(game.ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var mrtp, _ = flags.GetFloat64("reels")
				CalcStat(ctx, mrtp)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
