//go:build !prod || full || novomatic

package beetlemania

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Beetle Mania", Date: game.Year(2009)}, // see: https://www.slotsmate.com/software/novomatic/beetle-mania
		{Prov: "Novomatic", Name: "Beetle Mania Deluxe", Date: game.Date(2007, 11, 13)},
		{Prov: "Novomatic", Name: "Hot Target", Date: game.Date(2009, 9, 12)}, // see: https://www.slotsmate.com/software/novomatic/hot-target
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfghas |
			game.GPfgreel |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
}
