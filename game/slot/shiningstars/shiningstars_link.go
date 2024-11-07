//go:build !prod || full || agt

package shiningstars

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "agt/shiningstars", Prov: "AGT", Name: "Shining Stars"},
		{ID: "agt/greenhot", Prov: "AGT", Name: "Green Hot"},        // see: https://demo.agtsoftware.com/games/agt/greenhot
		{ID: "agt/applesshine", Prov: "AGT", Name: "Apples' Shine"}, // see: https://demo.agtsoftware.com/games/agt/applesshine
	},
	GP: game.GPfgno |
		game.GPscat |
		game.GPrwild,
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
