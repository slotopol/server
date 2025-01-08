//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dolphins Pearl"},        // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl
		{Prov: "Novomatic", Name: "Dolphins Pearl Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl-deluxe
		{Prov: "Aristocrat", Name: "Dolphin Treasure"},     // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
		{Prov: "Novomatic", Name: "Attila"},
		{Prov: "Novomatic", Name: "Banana Splash"},
		{Prov: "Novomatic", Name: "Dynasty Of Ming"},
		{Prov: "Novomatic", Name: "Gryphons Gold"},        // see: https://www.slotsmate.com/software/novomatic/gryphons-gold
		{Prov: "Novomatic", Name: "Gryphons Gold Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/gryphons-gold-deluxe
		{Prov: "Novomatic", Name: "Joker Dolphin"},
		{Prov: "Novomatic", Name: "King Of Cards"},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm"},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/lucky-ladys-charm-deluxe
		{Prov: "Novomatic", Name: "Pharaoh's Gold II"},
		{Prov: "Novomatic", Name: "Pharaoh's Gold III"},
		{Prov: "Novomatic", Name: "Polar Fox"},
		{Prov: "Novomatic", Name: "Ramses II"},
		{Prov: "Novomatic", Name: "Royal Treasures"},       // see: https://www.slotsmate.com/software/novomatic/novomatic-royal-treasures
		{Prov: "Novomatic", Name: "Secret Forest"},         // see: https://www.slotsmate.com/software/novomatic/secret-forest
		{Prov: "Novomatic", Name: "The Money Game"},        // see: https://www.slotsmate.com/software/novomatic/the-money-game
		{Prov: "Novomatic", Name: "The Money Game Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/the-money-game-deluxe
		{Prov: "Novomatic", Name: "Unicorn Magic"},
		{Prov: "Novomatic", Name: "Cold Spell"},      // see: https://www.slotsmate.com/software/novomatic/cold-spell
		{Prov: "Novomatic", Name: "Mermaid's Pearl"}, // see: https://www.slotsmate.com/software/novomatic/mermaids-pearl
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPfgreel |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
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
