//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "dolphinspearl", Prov: "Novomatic", Name: "Dolphins Pearl"},
		{ID: "dolphinspearldeluxe", Prov: "Novomatic", Name: "Dolphins Pearl Deluxe"},
		{ID: "dolphintreasure", Prov: "Aristocrat", Name: "Dolphin Treasure"}, // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
		{ID: "attila", Prov: "Novomatic", Name: "Attila"},
		{ID: "bananasplash", Prov: "Novomatic", Name: "Banana Splash"},
		{ID: "dynastyofming", Prov: "Novomatic", Name: "Dynasty Of Ming"},
		{ID: "gryphonsgold", Prov: "Novomatic", Name: "Gryphons Gold"},
		{ID: "jokerdolphin", Prov: "Novomatic", Name: "Joker Dolphin"},
		{ID: "kingofcards", Prov: "Novomatic", Name: "King Of Cards"},
		{ID: "luckyladyscharm", Prov: "Novomatic", Name: "Lucky Lady's Charm"},
		{ID: "luckyladyscharmdeluxe", Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe"},
		{ID: "pharaonsgold2", Prov: "Novomatic", Name: "Pharaon's Gold II"},
		{ID: "pharaonsgold3", Prov: "Novomatic", Name: "Pharaon's Gold III"},
		{ID: "polarfox", Prov: "Novomatic", Name: "Polar Fox"},
		{ID: "ramses2", Prov: "Novomatic", Name: "Ramses II"},
		{ID: "royaltreasures", Prov: "Novomatic", Name: "Royal Treasures"},
		{ID: "secretforest", Prov: "Novomatic", Name: "Secret Forest"},
		{ID: "themoneygame", Prov: "Novomatic", Name: "The Money Game"},
		{ID: "unicornmagic", Prov: "Novomatic", Name: "Unicorn Magic"},
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
