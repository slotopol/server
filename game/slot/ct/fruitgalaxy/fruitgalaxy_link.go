//go:build !prod || full || ct

package fruitgalaxy

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fruit Galaxy", Date: game.Date(2021, 5, 5)},   // see: https://www.slotsmate.com/software/ct-interactive/fruit-galaxy
		{Prov: "CT Interactive", Name: "Lord of Luck", Date: game.Date(2020, 12, 18)}, // see: https://www.slotsmate.com/software/ct-interactive/lord-of-luck
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX:  5,
		SY:  4,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/fruitgalaxy/reel"] = &ReelsMap
}
