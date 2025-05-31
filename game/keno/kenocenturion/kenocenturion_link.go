//go:build !prod || full || keno

package kenocenturion

import (
	"context"

	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Slotopol", Name: "Keno Centurion", Date: game.Year(2024)},
	},
	AlgDescr: game.AlgDescr{
		GT:  game.GTkeno,
		GP:  0,
		SX:  80,
		SY:  0,
		SN:  0,
		LN:  0,
		BN:  0,
		RTP: []float64{97.980099},
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, func(ctx context.Context, mrtp float64) float64 {
		return Paytable.CalcStat(ctx)
	})
}
