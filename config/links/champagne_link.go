//go:build !prod || full || megajack

package links

import (
	"context"

	"github.com/slotopol/server/game/champagne"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"champagne", "Champagne"},
		},
		Provider: "Megajack",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"96",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				champagne.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := champagne.ReelsMap[rd]; ok {
				return champagne.NewGame(rd)
			}
			return nil
		}
	}
}
