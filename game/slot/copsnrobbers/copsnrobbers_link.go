//go:build !prod || full || playngo

package copsnrobbers

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

func init() {
	var gi = game.GameInfo{
		Aliases: []game.GameAlias{
			{ID: "copsnrobbers", Name: "Cops'n'Robbers"},
		},
		Provider: "Play'n GO",
		SX:       5,
		SY:       3,
		GP: game.GPsel |
			game.GPfghas |
			game.GPfgreel |
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
				CalcStatReg(ctx, mrtp)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
