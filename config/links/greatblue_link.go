//go:build !prod || full || playtech

package links

import (
	"context"

	"github.com/slotopol/server/game/greatblue"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"greatblue", "Great Blue"}, // see: https://freeslotshub.com/playtech/great-blue/
			{"irishluck", "Irish Luck"}, // see: https://freeslotshub.com/playtech/irish-luck/
		},
		Provider: "Playtech",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"87", "89", "92", "94", "96", "97", "100", "108",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				greatblue.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return greatblue.NewGame(rtp)
		}
	}
}
