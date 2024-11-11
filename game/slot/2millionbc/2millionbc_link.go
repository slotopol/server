//go:build !prod || full || betsoft

package twomillionbc

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "BetSoft", Name: "2 Million B.C."},
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPscat,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  2,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		game.ScanFactory[aid] = CalcStatReg
		game.GameFactory[aid] = func() any { return NewGame() }
	}
}
