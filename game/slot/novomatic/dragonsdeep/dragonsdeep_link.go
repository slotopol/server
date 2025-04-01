//go:build !prod || full || novomatic

package dragonsdeep

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dragon's Deep"},
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPretrig |
		game.GPscat |
		game.GPwild |
		game.GPwturn,
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
