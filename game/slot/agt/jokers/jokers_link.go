//go:build !prod || full || agt

package jokers

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Jokers"},
		{Prov: "AGT", Name: "Happy Santa"}, // see: https://demo.agtsoftware.com/games/agt/happysanta
		{Prov: "AGT", Name: "Bigfoot"},     // see: https://demo.agtsoftware.com/games/agt/bigfoot
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPfgno |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
