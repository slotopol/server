//go:build !prod || full || playtech

package panthermoon

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Panther Moon", Date: game.Date(2009, 2, 28)}, // see: https://www.slotsmate.com/software/playtech/panther-moon
		{Prov: "Playtech", Name: "Safari Heat", Date: game.Date(2009, 3, 1)},   // see: https://www.slotsmate.com/software/playtech/safari-heat
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["playtech/panthermoon/bon"] = &ReelsBon
	game.DataRouter["playtech/panthermoon/reel"] = &ReelsMap
}
