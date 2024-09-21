//go:build !prod || full || netent

package links

import (
	"context"

	slot "github.com/slotopol/server/game/slot/trolls"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{ID: "trolls", Name: "Trolls"},
			{ID: "excalibur", Name: "Excalibur"},
			{ID: "pandorasbox", Name: "Pandora's Box"},
			{ID: "wildwitches", Name: "Wild Witches"},
		},
		Provider: "NetEnt",
		ScrnX:    5,
		ScrnY:    3,
		RtpList:  MakeRtpList(slot.ReelsMap),
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				slot.CalcStatReg(ctx, rn)
			}
		})
		GameFactory[ga.ID] = func() any {
			return slot.NewGame()
		}
	}
}
