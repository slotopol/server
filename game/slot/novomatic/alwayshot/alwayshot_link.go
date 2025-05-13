//go:build !prod || full || novomatic

package alwayshot

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Always Hot"},
		{Prov: "Novomatic", Name: "Always Hot Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/always-hot-deluxe
		{Prov: "Novomatic", Name: "Always American"},   // see: https://www.slotsmate.com/software/novomatic/always-american
		{Prov: "AGT", Name: "Tropic Hot"},              // see: https://demo.agtsoftware.com/games/agt/tropichot
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
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
