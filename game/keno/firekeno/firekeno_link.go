//go:build !prod || full || keno

package firekeno

import (
	"context"

	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Slotopol", Name: "Fire Keno"},
	},
	GP:  0,
	SX:  80,
	SY:  0,
	SN:  0,
	LN:  0,
	BN:  0,
	RTP: []float64{92.028857},
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, func(ctx context.Context, mrtp float64) float64 {
		return Paytable.CalcStat(ctx)
	})
}
