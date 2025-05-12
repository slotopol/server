//go:build !prod || full || netent

package diamonddogs

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Diamond Dogs", Year: 2013},
		{Prov: "NetEnt", Name: "Voodoo Vibes", Year: 2009},
	},
	GP: game.GPlpay |
		game.GPlsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPfgmult |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  1,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
