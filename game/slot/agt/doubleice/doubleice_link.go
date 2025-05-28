//go:build !prod || full || agt

package doubleice

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Double Ice"},
		{Prov: "AGT", Name: "Double Hot"}, // see: https://demo.agtsoftware.com/games/agt/double
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfill |
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
