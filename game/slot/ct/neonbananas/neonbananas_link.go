//go:build !prod || full || ct

package neonbananas

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Neon Bananas", Date: game.Date(2021, 6, 6)},  // see: https://www.slotsmate.com/software/ct-interactive/neon-bananas
		{Prov: "CT Interactive", Name: "Mighty Moon", Date: game.Date(2021, 7, 7)},   // see: https://www.slotsmate.com/software/ct-interactive/mighty-moon
		{Prov: "CT Interactive", Name: "Clover Wheel", Date: game.Date(2020, 12, 4)}, // see: https://www.slotsmate.com/software/ct-interactive/clover-wheel
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 2,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/neonbananas/reel"] = &ReelsMap
}
