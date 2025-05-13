//go:build !prod || full || netent

package trolls

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Trolls", Year: 2009},    // see: https://casino.ru/trolls-netent/
		{Prov: "NetEnt", Name: "Excalibur", Year: 2013}, // see: https://casino.ru/excalibur-netent/
		{Prov: "NetEnt", Name: "Pandora's Box", Year: 2009},
		{Prov: "NetEnt", Name: "Wild Witches", Year: 2010},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX:  5,
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
