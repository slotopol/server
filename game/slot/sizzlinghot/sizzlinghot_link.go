//go:build !prod || full || novomatic

package sizzlinghot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "sizzlinghot", Prov: "Novomatic", Name: "Sizzling Hot"},
		{ID: "sizzlinghotdeluxe", Prov: "Novomatic", Name: "Sizzling Hot Deluxe"},
	},
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
		game.ScanFactory[ga.ID] = CalcStat
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
