//go:build !prod || full || novomatic

package royaldynasty

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Royal Dynasty"}, // see: https://www.slotsmate.com/software/novomatic/royal-dynasty
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPscat |
		game.GPwild |
		game.GPwmult |
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
