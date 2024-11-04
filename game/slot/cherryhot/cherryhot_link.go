//go:build !prod || full || agt

package cherryhot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "cherryhot", Prov: "AGT", Name: "Cherry Hot"},
	},
	GP: game.GPsel |
		game.GPfgno |
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
