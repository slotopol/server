//go:build !prod || full || agt

package aislot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "agt/aislot", Prov: "AGT", Name: "AI"},   // see: https://demo.agtsoftware.com/games/agt/aislot
		{ID: "agt/tesla", Prov: "AGT", Name: "Tesla"}, // see: https://demo.agtsoftware.com/games/agt/tesla
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPscat |
		game.GPwild,
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
		game.ScanFactory[ga.ID] = CalcStat
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
