//go:build !prod || full || keno

package links

import (
	"context"

	keno "github.com/slotopol/server/game/keno/kenoluxury"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"kenoluxury", "Keno Luxury"},
			{"kenosports", "Keno Sports"},
		},
		Provider: "Slotopol",
		ScrnX:    80,
		ScrnY:    0,
		RtpList:  []float64{92.452462},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				keno.CalcStat(ctx)
			}
		})
		GameFactory[ga.ID] = func() any {
			return keno.NewGame()
		}
	}
}