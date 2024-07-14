//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/fortuneteller"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"fortuneteller", "Fortune Teller"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"85", "87", "89", "90", "91", "92", "94", "95", "96", "98", "99", "112", "126",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				fortuneteller.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			if _, ok := fortuneteller.ReelsMap[rd]; ok {
				return fortuneteller.NewGame(rd)
			}
			return nil
		}
	}
}
