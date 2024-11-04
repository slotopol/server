//go:build !prod || full || novomatic

package beetlemania

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "beetlemania", Prov: "Novomatic", Name: "Beetle Mania"},
		{ID: "beetlemaniadeluxe", Prov: "Novomatic", Name: "Beetle Mania Deluxe"},
		{ID: "hottarget", Prov: "Novomatic", Name: "Hot Target"},
	},
	GP: game.GPsel |
		game.GPfghas |
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
