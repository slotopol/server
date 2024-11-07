//go:build !prod || full || keno

package kenoluxury

import (
	"context"

	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "slotopol/kenoluxury", Prov: "Slotopol", Name: "Keno Luxury"},
		{ID: "slotopol/kenosports", Prov: "Slotopol", Name: "Keno Sports"},
	},
	GP:  0,
	SX:  80,
	SY:  0,
	SN:  0,
	LN:  0,
	BN:  0,
	RTP: []float64{92.104554},
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		game.ScanFactory[ga.ID] = func(ctx context.Context, mrtp float64) float64 {
			return Paytable.CalcStat(ctx)
		}
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
