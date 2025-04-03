//go:build !prod || full || netent

package fruitshop

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Fruit Shop"}, // see: https://www.slotsmate.com/software/netent/fruit-shop
	},
	GP: game.GPlpay |
		game.GPretrig |
		game.GPfgmult |
		game.GPwild |
		game.GPwmult,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
