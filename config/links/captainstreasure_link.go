//go:build !prod || full || novomatic

package links

import (
	"context"

	"github.com/slotopol/server/game/captainstreasure"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"captainstreasure", "Captain’s Treasure"},
		},
		Provider: "Playtech",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"91", "92", "93", "94", "95", "96", "97", "98", "99", "100", "112",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				captainstreasure.CalcStat(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := captainstreasure.ReelsMap[rd]; ok {
				return captainstreasure.NewGame(rd)
			}
			return nil
		}
	}
}
