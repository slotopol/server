//go:build !prod || full || agt

package jokers100

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Jokers"},
		{Prov: "AGT", Name: "50 Happy Santa"}, // see: https://demo.agtsoftware.com/games/agt/happysanta50
		{Prov: "AGT", Name: "40 Bigfoot"},     // see: https://demo.agtsoftware.com/games/agt/bigfoot40
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
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
