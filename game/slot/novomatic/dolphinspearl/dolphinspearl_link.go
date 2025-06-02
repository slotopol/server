//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dolphins Pearl", Date: game.Date(2001, 4, 25)},        // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl
		{Prov: "Novomatic", Name: "Dolphins Pearl Deluxe", Date: game.Date(2009, 4, 28)}, // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl-deluxe
		{Prov: "Novomatic", Name: "Attila", Date: game.Year(2009)},                       // see: https://casino.ru/attila-novomatic/
		{Prov: "Novomatic", Name: "Banana Splash", Date: game.Year(2009)},                // see: https://casino.ru/banana-splash-novomatic/
		{Prov: "Novomatic", Name: "Dynasty Of Ming", Date: game.Date(2008, 2, 20)},
		{Prov: "Novomatic", Name: "Gryphons Gold", Date: game.Year(2009)},              // see: https://www.slotsmate.com/software/novomatic/gryphons-gold
		{Prov: "Novomatic", Name: "Gryphons Gold Deluxe", Date: game.Date(2017, 4, 1)}, // see: https://www.slotsmate.com/software/novomatic/gryphons-gold-deluxe
		{Prov: "Novomatic", Name: "Joker Dolphin"},
		{Prov: "Novomatic", Name: "King Of Cards", Date: game.Date(2012, 2, 2)},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm", Date: game.Year(2001)},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe", Date: game.Year(2008)}, // see: https://www.slotsmate.com/software/novomatic/lucky-ladys-charm-deluxe
		{Prov: "Novomatic", Name: "Pharaoh's Gold II", Date: game.Year(2008)},         // see: https://casino.ru/pharaohs-gold-2-novomatic/
		{Prov: "Novomatic", Name: "Pharaoh's Gold III", Date: game.Year(2011)},
		{Prov: "Novomatic", Name: "Polar Fox", Date: game.Year(2008)}, // see: https://casino.ru/polar-fox-novomatic/
		{Prov: "Novomatic", Name: "Ramses II", Date: game.Date(2011, 4, 15)},
		{Prov: "Novomatic", Name: "Royal Treasures", Date: game.Year(2012)},             // see: https://www.slotsmate.com/software/novomatic/novomatic-royal-treasures
		{Prov: "Novomatic", Name: "Secret Forest", Date: game.Year(2013)},               // see: https://www.slotsmate.com/software/novomatic/secret-forest
		{Prov: "Novomatic", Name: "The Money Game", Date: game.Year(2009)},              // see: https://www.slotsmate.com/software/novomatic/the-money-game
		{Prov: "Novomatic", Name: "The Money Game Deluxe", Date: game.Date(2018, 7, 1)}, // see: https://www.slotsmate.com/software/novomatic/the-money-game-deluxe
		{Prov: "Novomatic", Name: "Unicorn Magic", Date: game.Year(2006)},               // see: https://casino.ru/unicorn-magic-novomatic/
		{Prov: "Novomatic", Name: "Cold Spell", Date: game.Date(2018, 10, 18)},          // see: https://www.slotsmate.com/software/novomatic/cold-spell
		{Prov: "Novomatic", Name: "Mermaid's Pearl", Date: game.Date(2014, 8, 15)},      // see: https://www.slotsmate.com/software/novomatic/mermaids-pearl
		{Prov: "Aristocrat", Name: "Dolphin Treasure", Date: game.Date(1996, 12, 1)},    // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
