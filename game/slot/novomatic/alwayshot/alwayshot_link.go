//go:build !prod || full || novomatic

package alwayshot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Always Hot"},
		{Prov: "AGT", Name: "Tropic Hot"}, // see: https://demo.agtsoftware.com/games/agt/tropichot
	},
	GP:  game.GPfgno,
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
