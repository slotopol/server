//go:build !prod || full || novomatic

package alwayshot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "alwayshot", Name: "Always Hot"},
	},
	Provider: "Novomatic",
	GP:       game.GPfgno,
	SX:       3,
	SY:       3,
	SN:       len(LinePay),
	LN:       len(BetLines),
	BN:       0,
	RTP:      game.MakeRtpList(ReelsMap),
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		game.ScanFactory[ga.ID] = CalcStat
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
