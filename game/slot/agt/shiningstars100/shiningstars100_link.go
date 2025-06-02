//go:build !prod || full || agt

package shiningstars100

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Shining Stars"},
		{Prov: "AGT", Name: "50 Apples' Shine"},                 // see: https://demo.agtsoftware.com/games/agt/applesshine50
		{Prov: "AGT", Name: "Red Crown", Date: game.Year(2025)}, // see: https://demo.agtsoftware.com/games/agt/redcrown
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
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
}
