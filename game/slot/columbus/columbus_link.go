//go:build !prod || full || novomatic

package columbus

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "novomatic/columbus", Prov: "Novomatic", Name: "Columbus"},
		{ID: "novomatic/columbusdeluxe", Prov: "Novomatic", Name: "Columbus Deluxe"},
		{ID: "novomatic/marcopolo", Prov: "Novomatic", Name: "Marco Polo"},
	},
	GP: game.GPsel |
		game.GPretrig |
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
