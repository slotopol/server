//go:build !prod || full || novomatic

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/dolphinspearl"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
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
		GP:       GPsel | GPretrig | GPfgmult | GPfgreel | GPscat | GPwild,
		SN:       len(slot.LinePay),
		LN:       10,
		BN:       0,
		RTP:      MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					slot.CalcStatBon(ctx)
				} else {
					slot.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
