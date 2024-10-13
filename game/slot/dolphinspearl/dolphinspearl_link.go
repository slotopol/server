//go:build !prod || full || novomatic

package dolphinspearl

import (
	"context"

	game "github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

func init() {
	var gi = game.GameInfo{
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
		SX:       5,
		SY:       3,
		GP: game.GPsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPfgreel |
			game.GPscat |
			game.GPwild,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	}
	game.GameList = append(game.GameList, gi)

	for _, ga := range gi.Aliases {
		game.ScanIters = append(game.ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var mrtp, _ = flags.GetFloat64("reels")
				CalcStatReg(ctx, mrtp)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
