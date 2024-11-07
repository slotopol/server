//go:build !prod || full || playtech

package greatblue

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "playtech/greatblue", Prov: "Playtech", Name: "Great Blue"}, // see: https://freeslotshub.com/playtech/great-blue/
		{ID: "playtech/irishluck", Prov: "Playtech", Name: "Irish Luck"}, // see: https://freeslotshub.com/playtech/irish-luck/
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
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
		game.ScanFactory[ga.ID] = CalcStat
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
