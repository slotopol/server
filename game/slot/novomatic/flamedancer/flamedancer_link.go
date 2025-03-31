//go:build !prod || full || novomatic

package flamedancer

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Flame Dancer"}, // see: https://casino.ru/flame-dancer-novomatic/
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPscat |
		game.GPwild |
		game.GPrwild |
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
