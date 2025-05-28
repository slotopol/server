//go:build !prod || full || agt

package aladdin

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Aladdin"},          // see: https://demo.agtsoftware.com/games/agt/aladdin
		{Prov: "AGT", Name: "Wild West"},        // see: https://demo.agtsoftware.com/games/agt/wildwest
		{Prov: "AGT", Name: "Crown"},            // see: https://demo.agtsoftware.com/games/agt/crown
		{Prov: "AGT", Name: "Arabian Nights 2"}, // see: https://demo.agtsoftware.com/games/agt/arabiannights2
		{Prov: "AGT", Name: "Casino"},           // see: https://demo.agtsoftware.com/games/agt/casino
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
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
