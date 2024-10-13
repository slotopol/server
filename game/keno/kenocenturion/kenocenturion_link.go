//go:build !prod || full || keno

package kenocenturion

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

func init() {
	var gi = game.GameInfo{
		Aliases: []game.GameAlias{
			{ID: "kenocenturion", Name: "Keno Centurion"},
		},
		Provider: "Slotopol",
		SX:       80,
		SY:       0,
		GP:       0,
		SN:       0,
		LN:       0,
		BN:       0,
		RTP:      []float64{97.980099},
	}
	game.GameList = append(game.GameList, gi)

	for _, ga := range gi.Aliases {
		game.ScanIters = append(game.ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				CalcStat(ctx)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
