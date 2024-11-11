//go:build !prod || full || agt

package shiningstars100

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Shining Stars"},
		{Prov: "AGT", Name: "50 Apples' Shine"}, // see: https://demo.agtsoftware.com/games/agt/applesshine50
	},
	GP: game.GPsel |
		game.GPfgno |
		game.GPscat |
		game.GPrwild,
	SX:  5,
	SY:  4,
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
