//go:build !prod || full || novomatic

package sizzlinghot

import (
	"context"

	"github.com/slotopol/server/game"
	"github.com/spf13/pflag"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "sizzlinghot", Name: "Sizzling Hot"},
		{ID: "sizzlinghotdeluxe", Name: "Sizzling Hot Deluxe"},
	},
	Provider: "Novomatic",
	GP: game.GPfgno |
		game.GPscat,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	game.GameList = append(game.GameList, &Info)

	for _, ga := range Info.Aliases {
		game.ScanIters = append(game.ScanIters, func(flags *pflag.FlagSet, ctx context.Context) {
			if is, _ := flags.GetBool(ga.ID); is {
				var mrtp, _ = flags.GetFloat64("reels")
				CalcStat(ctx, mrtp)
			}
		})
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
