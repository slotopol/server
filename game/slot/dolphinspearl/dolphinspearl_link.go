//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "novomatic/dolphinspearl", Prov: "Novomatic", Name: "Dolphins Pearl"},
		{ID: "novomatic/dolphinspearldeluxe", Prov: "Novomatic", Name: "Dolphins Pearl Deluxe"},
		{ID: "aristocrat/dolphintreasure", Prov: "Aristocrat", Name: "Dolphin Treasure"}, // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
		{ID: "novomatic/attila", Prov: "Novomatic", Name: "Attila"},
		{ID: "novomatic/bananasplash", Prov: "Novomatic", Name: "Banana Splash"},
		{ID: "novomatic/dynastyofming", Prov: "Novomatic", Name: "Dynasty Of Ming"},
		{ID: "novomatic/gryphonsgold", Prov: "Novomatic", Name: "Gryphons Gold"},
		{ID: "novomatic/jokerdolphin", Prov: "Novomatic", Name: "Joker Dolphin"},
		{ID: "novomatic/kingofcards", Prov: "Novomatic", Name: "King Of Cards"},
		{ID: "novomatic/luckyladyscharm", Prov: "Novomatic", Name: "Lucky Lady's Charm"},
		{ID: "novomatic/luckyladyscharmdeluxe", Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe"},
		{ID: "novomatic/pharaonsgold2", Prov: "Novomatic", Name: "Pharaon's Gold II"},
		{ID: "novomatic/pharaonsgold3", Prov: "Novomatic", Name: "Pharaon's Gold III"},
		{ID: "novomatic/polarfox", Prov: "Novomatic", Name: "Polar Fox"},
		{ID: "novomatic/ramses2", Prov: "Novomatic", Name: "Ramses II"},
		{ID: "novomatic/royaltreasures", Prov: "Novomatic", Name: "Royal Treasures"},
		{ID: "novomatic/secretforest", Prov: "Novomatic", Name: "Secret Forest"},
		{ID: "novomatic/themoneygame", Prov: "Novomatic", Name: "The Money Game"},
		{ID: "novomatic/unicornmagic", Prov: "Novomatic", Name: "Unicorn Magic"},
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
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
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		game.ScanFactory[ga.ID] = CalcStatReg
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
