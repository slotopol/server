//go:build !prod || full || keno

package kenocenturion

import (
	"context"

	"github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "kenocenturion", Name: "Keno Centurion"},
	},
	Provider: "Slotopol",
	GP:       0,
	SX:       80,
	SY:       0,
	SN:       0,
	LN:       0,
	BN:       0,
	RTP:      []float64{97.980099},
}

func init() {
	game.GameList = append(game.GameList, &Info)

	for _, ga := range Info.Aliases {
		game.ScanIters = append(game.ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				Paytable.CalcStat(ctx)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
