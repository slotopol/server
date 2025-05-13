//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dolphins Pearl", Year: 2001},        // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl
		{Prov: "Novomatic", Name: "Dolphins Pearl Deluxe", Year: 2006}, // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl-deluxe
		{Prov: "Novomatic", Name: "Attila", Year: 2009},                // see: https://casino.ru/attila-novomatic/
		{Prov: "Novomatic", Name: "Banana Splash", Year: 2009},         // see: https://casino.ru/banana-splash-novomatic/
		{Prov: "Novomatic", Name: "Dynasty Of Ming", Year: 2008},
		{Prov: "Novomatic", Name: "Gryphons Gold", Year: 2009}, // see: https://www.slotsmate.com/software/novomatic/gryphons-gold
		{Prov: "Novomatic", Name: "Gryphons Gold Deluxe"},      // see: https://www.slotsmate.com/software/novomatic/gryphons-gold-deluxe
		{Prov: "Novomatic", Name: "Joker Dolphin"},
		{Prov: "Novomatic", Name: "King Of Cards"},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm", Year: 2001},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe", Year: 2008}, // see: https://www.slotsmate.com/software/novomatic/lucky-ladys-charm-deluxe
		{Prov: "Novomatic", Name: "Pharaoh's Gold II", Year: 2008},         // see: https://casino.ru/pharaohs-gold-2-novomatic/
		{Prov: "Novomatic", Name: "Pharaoh's Gold III"},
		{Prov: "Novomatic", Name: "Polar Fox", Year: 2008}, // see: https://casino.ru/polar-fox-novomatic/
		{Prov: "Novomatic", Name: "Ramses II"},
		{Prov: "Novomatic", Name: "Royal Treasures"},                   // see: https://www.slotsmate.com/software/novomatic/novomatic-royal-treasures
		{Prov: "Novomatic", Name: "Secret Forest"},                     // see: https://www.slotsmate.com/software/novomatic/secret-forest
		{Prov: "Novomatic", Name: "The Money Game"},                    // see: https://www.slotsmate.com/software/novomatic/the-money-game
		{Prov: "Novomatic", Name: "The Money Game Deluxe", Year: 2012}, // see: https://www.slotsmate.com/software/novomatic/the-money-game-deluxe
		{Prov: "Novomatic", Name: "Unicorn Magic", Year: 2006},         // see: https://casino.ru/unicorn-magic-novomatic/
		{Prov: "Novomatic", Name: "Cold Spell"},                        // see: https://www.slotsmate.com/software/novomatic/cold-spell
		{Prov: "Novomatic", Name: "Mermaid's Pearl"},                   // see: https://www.slotsmate.com/software/novomatic/mermaids-pearl
		{Prov: "Aristocrat", Name: "Dolphin Treasure"},                 // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
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
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
