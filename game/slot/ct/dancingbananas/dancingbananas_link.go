//go:build !prod || full || ct

package dancingbananas

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Dancing Bananas", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/dancing-bananas
		{Prov: "CT Interactive", Name: "Jaguar Warrior", Date: game.Date(2020, 11, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/jaguar-warrior
		{Prov: "CT Interactive", Name: "Clover Joker", Date: game.Date(2021, 8, 8)},      // see: https://www.slotsmate.com/software/ct-interactive/clover-joker
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/dancingbananas/reel"] = &ReelsMap
}
