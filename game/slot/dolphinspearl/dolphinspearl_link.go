//go:build !prod || full || novomatic

package dolphinspearl

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "dolphinspearl", Name: "Dolphins Pearl"},
		{ID: "dolphinspearldeluxe", Name: "Dolphins Pearl Deluxe"},
		{ID: "dolphintreasure", Name: "Dolphin Treasure"}, // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
		{ID: "attila", Name: "Attila"},
		{ID: "bananasplash", Name: "Banana Splash"},
		{ID: "dynastyofming", Name: "Dynasty Of Ming"},
		{ID: "gryphonsgold", Name: "Gryphons Gold"},
		{ID: "jokerdolphin", Name: "Joker Dolphin"},
		{ID: "kingofcards", Name: "King Of Cards"},
		{ID: "luckyladyscharm", Name: "Lucky Lady's Charm"},
		{ID: "luckyladyscharmdeluxe", Name: "Lucky Lady's Charm Deluxe"},
		{ID: "pharaonsgold2", Name: "Pharaon's Gold II"},
		{ID: "pharaonsgold3", Name: "Pharaon's Gold III"},
		{ID: "polarfox", Name: "Polar Fox"},
		{ID: "ramses2", Name: "Ramses II"},
		{ID: "royaltreasures", Name: "Royal Treasures"},
		{ID: "secretforest", Name: "Secret Forest"},
		{ID: "themoneygame", Name: "The Money Game"},
		{ID: "unicornmagic", Name: "Unicorn Magic"},
	},
	Provider: "Novomatic",
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
