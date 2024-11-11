//go:build !prod || full || keno

package kenocenturion

import (
	"context"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Slotopol", Name: "Keno Centurion"},
	},
	GP:  0,
	SX:  80,
	SY:  0,
	SN:  0,
	LN:  0,
	BN:  0,
	RTP: []float64{97.980099},
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		game.ScanFactory[aid] = func(ctx context.Context, mrtp float64) float64 {
			return Paytable.CalcStat(ctx)
		}
		game.GameFactory[aid] = func() any { return NewGame() }
	}
}
