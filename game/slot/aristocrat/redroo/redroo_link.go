//go:build !prod || full || aristocrat

package redroo

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Aristocrat", Name: "Redroo"},
	},
	GP: game.GPretrig |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
	SX:  5,
	SY:  4,
	SN:  len(LinePay),
	LN:  1024,
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
