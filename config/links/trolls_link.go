//go:build !prod || full || netent

package links

import (
	"context"

	"github.com/slotopol/server/game/trolls"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"trolls", "Trolls"},
			{"excalibur", "Excalibur"},
			{"pandorasbox", "Pandora's Box"},
			{"wildwitches", "Wild Witches"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"88", "89", "92", "93", "94", "95", "97", "98", "102", "110",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				trolls.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func(rd string) any {
			return trolls.NewGame(Atof(rd))
		}
	}
}
