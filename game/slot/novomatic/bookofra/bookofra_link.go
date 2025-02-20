//go:build !prod || full || novomatic

package bookofra

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Book of Ra"},        // see: https://www.slotsmate.com/software/novomatic/book-of-ra-classic
		{Prov: "Novomatic", Name: "Book of Ra Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe
		{Prov: "Novomatic", Name: "Gate of Ra Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/gate-of-ra-deluxe
		{Prov: "Novomatic", Name: "Golden Prophecies"}, // see: https://www.slotsmate.com/software/novomatic/golden-prophecies
		{Prov: "Novomatic", Name: "Down Under"},
		{Prov: "Novomatic", Name: "God of Sun"},
		{Prov: "Novomatic", Name: "Lord of the Ocean"},
		{Prov: "Novomatic", Name: "Faust"},                      // see: https://freeslotshub.com/novomatic/faust/
		{Prov: "Novomatic", Name: "The Real King Gold Records"}, // see: https://www.slotsmate.com/software/novomatic/the-real-king-gold-records
		{Prov: "Novomatic", Name: "Angry Birds"},                // see: https://www.slotsmate.com/software/novomatic/angry-birds
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
