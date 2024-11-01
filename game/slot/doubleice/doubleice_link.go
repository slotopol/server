//go:build !prod || full || agt

package doubleice

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "doubleice", Name: "Double Ice"},
		{ID: "doublehot", Name: "Double Hot"}, // see: https://demo.agtsoftware.com/games/agt/double
	},
	Provider: "AGT",
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
