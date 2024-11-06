//go:build !prod || full || agt

package iceqween

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "icequeen", Prov: "AGT", Name: "Ice Queen"}, // see: https://demo.agtsoftware.com/games/agt/iceqween
		{ID: "stalker", Prov: "AGT", Name: "STALKER"},    // see: https://demo.agtsoftware.com/games/agt/stalker
		{ID: "bigfive", Prov: "AGT", Name: "Big Five"},   // see: https://demo.agtsoftware.com/games/agt/bigfive
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPfgreel |
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
		game.ScanFactory[ga.ID] = CalcStatReg
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
