//go:build !prod || full || keno

package kenofast

import (
	"context"

	"github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "kenofast", Name: "Keno Fast"},
	},
	Provider: "AGT",
	GP:       0,
	SX:       80,
	SY:       0,
	SN:       0,
	LN:       0,
	BN:       0,
	RTP:      []float64{95.616967},
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
