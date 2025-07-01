//go:build !prod || full || ct

package winstorm

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Win Storm", Date: game.Date(2020, 12, 27)}, // see: https://www.slotsmate.com/software/ct-interactive/win-storm
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcasc |
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
