//go:build !prod || full || netent

package groovysixties

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Groovy Sixties"},
		{Prov: "NetEnt", Name: "Funky Seventies"}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
		{Prov: "NetEnt", Name: "Super Eighties"},  // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
	},
	GP: game.GPlpay |
		game.GPsel |
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
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
