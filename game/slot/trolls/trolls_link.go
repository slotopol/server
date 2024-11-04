//go:build !prod || full || netent

package trolls

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "trolls", Prov: "NetEnt", Name: "Trolls"},
		{ID: "excalibur", Prov: "NetEnt", Name: "Excalibur"},
		{ID: "pandorasbox", Prov: "NetEnt", Name: "Pandora's Box"},
		{ID: "wildwitches", Prov: "NetEnt", Name: "Wild Witches"},
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
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
