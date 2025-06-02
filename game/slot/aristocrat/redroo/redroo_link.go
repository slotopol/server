//go:build !prod || full || aristocrat

package redroo

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Aristocrat", Name: "Redroo", Date: game.Date(2017, 6, 28)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPretrig |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX:  5,
		SY:  4,
		SN:  len(LinePay),
		WN:  1024,
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
}
