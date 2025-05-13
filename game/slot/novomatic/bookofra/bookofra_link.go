//go:build !prod || full || novomatic

package bookofra

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Book of Ra"},                    // see: https://www.slotsmate.com/software/novomatic/book-of-ra-classic
		{Prov: "Novomatic", Name: "Book of Ra Deluxe", Year: 2008}, // see: https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe
		{Prov: "Novomatic", Name: "Gate of Ra Deluxe"},             // see: https://www.slotsmate.com/software/novomatic/gate-of-ra-deluxe
		{Prov: "Novomatic", Name: "Golden Prophecies"},             // see: https://www.slotsmate.com/software/novomatic/golden-prophecies
		{Prov: "Novomatic", Name: "Down Under"},
		{Prov: "Novomatic", Name: "God of Sun"},
		{Prov: "Novomatic", Name: "Lord of the Ocean", Year: 2008}, // see: https://casino.ru/lord-of-the-ocean-novomatic/
		{Prov: "Novomatic", Name: "Faust", Year: 2009},             // see: https://freeslotshub.com/novomatic/faust/
		{Prov: "Novomatic", Name: "The Real King Gold Records"},    // see: https://www.slotsmate.com/software/novomatic/the-real-king-gold-records
		{Prov: "Novomatic", Name: "Angry Birds"},                   // see: https://www.slotsmate.com/software/novomatic/angry-birds
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
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
	},
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
