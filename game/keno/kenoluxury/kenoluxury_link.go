//go:build !prod || full || keno

package kenoluxury

import (
	"context"

	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Aristocrat", Name: "Keno Luxury", Date: game.Year(2013)},
		{Prov: "Aristocrat", Name: "Keno Sports", Date: game.Year(2013)},
	},
	AlgDescr: game.AlgDescr{
		GT:  game.GTkeno,
		GP:  0,
		SX:  80,
		SY:  0,
		SN:  0,
		LN:  0,
		BN:  0,
		RTP: []float64{92.104554},
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, func(ctx context.Context, mrtp float64) float64 {
		return Paytable.CalcStat(ctx)
	})
}
