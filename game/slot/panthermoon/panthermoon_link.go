//go:build !prod || full || playtech

package panthermoon

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "panthermoon", Prov: "Playtech", Name: "Panther Moon"},
		{ID: "safariheat", Prov: "Playtech", Name: "Safari Heat"},
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
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
