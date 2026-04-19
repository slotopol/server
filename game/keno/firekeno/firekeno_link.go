//go:build !prod || full || keno

package firekeno

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Slotopol", Name: "Fire Keno", Date: game.Year(2024)},
	},
	AlgDescr: game.AlgDescr{
		GT:  game.GTkeno,
		GP:  0,
		SX:  80,
		SY:  0,
		SN:  0,
		LN:  0,
		BN:  0,
		RTP: []float64{92.013465},
	},
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, Paytable.CalcStat)
}
