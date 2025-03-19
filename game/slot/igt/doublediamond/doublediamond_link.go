//go:build !prod || full || igt

package doublediamond

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Double Diamond"}, // see: https://www.slotsmate.com/software/igt/double-diamond
	},
	GP: game.GPfgno |
		game.GPwild |
		game.GPwmult,
	SX:  3,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
