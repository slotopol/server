//go:build !prod || full || novomatic

package justjewels

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Just Jewels"},        // see: https://www.slotsmate.com/software/novomatic/just-jewels
		{Prov: "Novomatic", Name: "Just Jewels Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/just-jewels-deluxe
		{Prov: "Novomatic", Name: "Just Fruits"},        // see: https://www.slotsmate.com/software/novomatic/just-fruits
	},
	GP: game.GPsel |
		game.GPfgno |
		game.GPscat,
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
