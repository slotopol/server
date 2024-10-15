//go:build !prod || full || keno || aristocrat

package americankeno

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "americankeno", Name: "American Keno"},
	},
	Provider: "Aristocrat",
	GP:       0,
	SX:       80,
	SY:       0,
	SN:       0,
	LN:       0,
	BN:       0,
	RTP:      []float64{89.250235},
}

func init() {
	game.GameList = append(game.GameList, &Info)

	for _, ga := range Info.Aliases {
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
