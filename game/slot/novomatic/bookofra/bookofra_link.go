//go:build !prod || full || novomatic

package bookofra

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Book Of Ra"},
		{Prov: "Novomatic", Name: "Book Of Ra Deluxe"}, // see: https://freeslotshub.com/novomatic/book-of-ra-deluxe/
		{Prov: "Novomatic", Name: "Down Under"},
		{Prov: "Novomatic", Name: "God Of Sun"},
		{Prov: "Novomatic", Name: "Lord of the Ocean"},
		{Prov: "Novomatic", Name: "Faust"}, // see: https://freeslotshub.com/novomatic/faust/
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPscat |
		game.GPwild,
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