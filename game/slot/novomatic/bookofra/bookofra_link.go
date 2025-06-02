//go:build !prod || full || novomatic

package bookofra

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Book of Ra", Date: game.Year(2005)},                // see: https://www.slotsmate.com/software/novomatic/book-of-ra-classic
		{Prov: "Novomatic", Name: "Book of Ra Deluxe", Date: game.Date(2008, 4, 11)},  // see: https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe
		{Prov: "Novomatic", Name: "Gate of Ra Deluxe", Date: game.Date(2018, 12, 15)}, // see: https://www.slotsmate.com/software/novomatic/gate-of-ra-deluxe
		{Prov: "Novomatic", Name: "Golden Prophecies", Date: game.Date(2018, 8, 15)},  // see: https://www.slotsmate.com/software/novomatic/golden-prophecies
		{Prov: "Novomatic", Name: "Down Under"},
		{Prov: "Novomatic", Name: "God of Sun"},
		{Prov: "Novomatic", Name: "Lord of the Ocean", Date: game.Date(2008, 6, 15)},          // see: https://casino.ru/lord-of-the-ocean-novomatic/
		{Prov: "Novomatic", Name: "Faust", Date: game.Date(2015, 2, 15)},                      // see: https://freeslotshub.com/novomatic/faust/
		{Prov: "Novomatic", Name: "The Real King Gold Records", Date: game.Date(2020, 2, 11)}, // see: https://www.slotsmate.com/software/novomatic/the-real-king-gold-records
		{Prov: "Novomatic", Name: "Angry Birds", Date: game.Date(2019, 3, 16)},                // see: https://www.slotsmate.com/software/novomatic/angry-birds
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
}
