//go:build !prod || full || novomatic

package plentyofjewels20

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Plenty of Jewels 20 hot", Date: game.Date(2015, 11, 1)}, // see: https://www.slotsmate.com/software/novomatic/plenty-of-jewels-20-hot
		{Prov: "Novomatic", Name: "Plenty of Fruit 20 hot", Date: game.Date(2015, 11, 4)},  // see: https://www.slotsmate.com/software/novomatic/plenty-of-fruit-20-hot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX:  5,
		SY:  3,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
}
