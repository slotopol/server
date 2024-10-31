//go:build !prod || full || agt

package iceiceice

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "iceiceice", Name: "Ice Ice Ice"},
		{ID: "5hothothot", Name: "5 Hot Hot Hot"}, // see: https://demo.agtsoftware.com/games/agt/hothothot5
	},
	Provider: "AGT",
	GP: game.GPretrig |
		game.GPscat |
		game.GPwild,
	SX:  3,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		game.ScanFactory[ga.ID] = CalcStatReg
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
