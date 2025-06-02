//go:build !prod || full || netent

package groovysixties

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Groovy Sixties", Date: game.Date(2009, 6, 15)},
		{Prov: "NetEnt", Name: "Funky Seventies", Date: game.Date(2009, 6, 15)}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
		{Prov: "NetEnt", Name: "Super Eighties", Date: game.Date(2009, 6, 15)},  // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX:  5,
		SY:  4,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
}
