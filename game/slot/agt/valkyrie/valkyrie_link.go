//go:build !prod || full || agt

package valkyrie

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Valkyrie", Date: game.Year(2024)}, // see: https://agtsoftware.com/games/agt/valkyrie
		{Prov: "AGT", Name: "Aquaman", Date: game.Year(2025)},  // see: https://agtsoftware.com/games/agt/aquaman
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfghas |
			game.GPscat |
			game.GPwild |
			game.GPbsym,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["agt/valkyrie/bon"] = &ReelsBon
	game.DataRouter["agt/valkyrie/reel"] = &ReelsMap
}
