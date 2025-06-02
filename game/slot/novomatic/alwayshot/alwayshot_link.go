//go:build !prod || full || novomatic

package alwayshot

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Always Hot", Date: game.Date(2008, 11, 17)},
		{Prov: "Novomatic", Name: "Always Hot Deluxe", Date: game.Date(2008, 11, 18)}, // see: https://www.slotsmate.com/software/novomatic/always-hot-deluxe
		{Prov: "Novomatic", Name: "Always American", Date: game.Date(2016, 4, 15)},    // see: https://www.slotsmate.com/software/novomatic/always-american
		{Prov: "AGT", Name: "Tropic Hot"},                                             // see: https://demo.agtsoftware.com/games/agt/tropichot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlsel |
			game.GPfgno,
		SX:  3,
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
