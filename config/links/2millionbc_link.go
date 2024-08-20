//go:build !prod || full || betsoft

package links

import (
	"context"

	twomillionbc "github.com/slotopol/server/game/2millionbc"
	"github.com/spf13/pflag"
)

func init() {
	var gi = GameInfo{
		Aliases: []GameAlias{
			{"2millionbc", "2 Million B.C."},
		},
		Provider: "BetSoft",
		ScrnX:    5,
		ScrnY:    3,
		RtpList: []string{
			"91", "93", "94", "96", "97", "100", "114", "bon",
		},
	}
	GameList = append(GameList, gi)

	for _, ga := range gi.Aliases {
		ScanIters = append(ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var rn, _ = flags.GetString("reels")
				if rn == "bon" {
					twomillionbc.CalcStatBon(ctx)
				} else {
					twomillionbc.CalcStatReg(ctx, rn)
				}
			}
		})
		GameFactory[ga.ID] = func(rtp float64) any {
			return twomillionbc.NewGame(rtp)
		}
	}
}
