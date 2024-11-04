//go:build !prod || full || novomatic

package roaringforties

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "roaringforties", Prov: "Novomatic", Name: "Roaring Forties"},
	},
	GP: game.GPsel |
		game.GPfgno |
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
		game.ScanFactory[ga.ID] = CalcStat
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
