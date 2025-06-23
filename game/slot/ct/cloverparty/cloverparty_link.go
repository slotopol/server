//go:build !prod || full || ct

package cloverparty

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Clover Party", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/clover-party
		{Prov: "CT Interactive", Name: "20 Mega Fresh", Date: game.Date(2021, 7, 7)},  // see: https://www.slotsmate.com/software/ct-interactive/20-mega-fresh
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
