//go:build !prod || full || keno

package links

import (
	"context"

	keno "github.com/slotopol/server/game/keno/kenocenturion"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "kenocenturion", Name: "Keno Centurion"},
		},
		Provider: "Slotopol",
		SX:       80,
		SY:       0,
		LN:       0,
		FG:       FGno,
		BN:       0,
		RTP:      []float64{97.980099},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				keno.CalcStat(ctx)
			}
		})
		GameFactory[ga.ID] = func() any {
			return keno.NewGame()
		}
	}
}
