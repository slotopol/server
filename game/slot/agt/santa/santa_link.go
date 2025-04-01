//go:build !prod || full || agt

package santa

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Santa"},
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPretrig |
		game.GPscat |
		game.GPwild,
	SX:  4,
	SY:  4,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
