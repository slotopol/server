//go:build !prod || full || novomatic

package chicago

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

func init() {
	var gi = game.GameInfo{
		Aliases: []game.GameAlias{
			{ID: "chicago", Name: "Chicago"},
		},
		Provider: "Novomatic",
		SX:       5,
		SY:       3,
		GP:       game.GPsel | 
		game.GPretrig | 
		game.GPfgmult | 
		game.GPscat | 
		game.GPwild,
		SN:       len(LinePay),
		LN:       len(BetLines),
		BN:       0,
		RTP:      game.MakeRtpList(ReelsMap),
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
