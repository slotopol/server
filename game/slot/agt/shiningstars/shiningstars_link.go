//go:build !prod || full || agt

package shiningstars

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Shining Stars"},
		{Prov: "AGT", Name: "Green Hot"},     // see: https://demo.agtsoftware.com/games/agt/greenhot
		{Prov: "AGT", Name: "Apples' Shine"}, // see: https://demo.agtsoftware.com/games/agt/applesshine
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
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		game.ScanFactory[aid] = CalcStat
		game.GameFactory[aid] = func() any { return NewGame() }
	}
}
