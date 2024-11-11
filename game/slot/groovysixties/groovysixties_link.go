//go:build !prod || full || netent

package groovysixties

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Groovy Sixties"},
		{Prov: "NetEnt", Name: "Funky Seventies"}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
		{Prov: "NetEnt", Name: "Super Eighties"},  // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPscat |
		game.GPwild,
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
