//go:build !prod || full || megajack

package slotopoldeluxe

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Slotopol Deluxe"},
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPjack |
		game.GPfgno |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  4,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
