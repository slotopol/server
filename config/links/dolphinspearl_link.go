//go:build !prod || full || novomatic

package links

import (
	"context"
	"strconv"

	"github.com/slotopol/server/game/dolphinspearl"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"dolphinspearl", "Dolphins Pearl"},
			{"dolphinspearldeluxe", "Dolphins Pearl Deluxe"},
			{"dolphintreasure", "Dolphin Treasure"}, // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
			{"attila", "Attila"},
			{"bananasplash", "Banana Splash"},
			{"dynastyofming", "Dynasty Of Ming"},
			{"gryphonsgold", "Gryphons Gold"},
			{"jokerdolphin", "Joker Dolphin"},
			{"kingofcards", "King Of Cards"},
			{"luckyladyscharm", "Lucky Lady's Charm"},
			{"luckyladyscharmdeluxe", "Lucky Lady's Charm Deluxe"},
			{"pharaonsgold2", "Pharaon's Gold II"},
			{"pharaonsgold3", "Pharaon's Gold III"},
			{"polarfox", "Polar Fox"},
			{"ramses2", "Ramses II"},
			{"royaltreasures", "Royal Treasures"},
			{"secretforest", "Secret Forest"},
			{"themoneygame", "The Money Game"},
			{"unicornmagic", "Unicorn Magic"},
		},
		Provider: "Novomatic",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"86", "88", "90", "92", "94", "95", "96", "97", "141", "bon",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					dolphinspearl.CalcStatBon(ctx)
				} else {
					dolphinspearl.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			var rtp, _ = strconv.ParseFloat(rd, 64)
			return dolphinspearl.NewGame(rtp)
		}
	}
}
