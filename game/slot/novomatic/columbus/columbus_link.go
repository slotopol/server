//go:build !prod || full || novomatic

package columbus

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Columbus", Date: game.Year(2005)},                        // see: https://casino.ru/columbus-novomatic/
		{Prov: "Novomatic", Name: "Columbus Deluxe", Date: game.Date(2008, 3, 19)},          // see: https://www.slotsmate.com/software/novomatic/columbus-deluxe
		{Prov: "Novomatic", Name: "Marco Polo", Date: game.Year(2008)},                      // see: https://casino.ru/marco-polo-novomatic/
		{Prov: "Novomatic", Name: "Holmes and Watson Deluxe", Date: game.Date(2018, 3, 15)}, // see: https://www.slotsmate.com/software/novomatic/holmes-and-watson-deluxe
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
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
