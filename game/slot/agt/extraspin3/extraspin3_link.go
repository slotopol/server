//go:build !prod || full || agt

package extraspin3

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slot/agt/extraspin"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Extra Spin III", LNum: 10}, // see: https://agtsoftware.com/games/agt/extraspin3
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(extraspin.ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, extraspin.CalcStat)
}
