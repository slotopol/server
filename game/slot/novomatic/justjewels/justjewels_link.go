//go:build !prod || full || novomatic

package justjewels

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Just Jewels", Date: game.Year(2005)},              // see: https://www.slotsmate.com/software/novomatic/just-jewels
		{Prov: "Novomatic", Name: "Just Jewels Deluxe", Date: game.Date(2008, 7, 1)}, // see: https://www.slotsmate.com/software/novomatic/just-jewels-deluxe
		{Prov: "Novomatic", Name: "Just Fruits", Date: game.Year(2001)},              // see: https://www.slotsmate.com/software/novomatic/just-fruits
		{Prov: "Novomatic", Name: "Royal Jewels"},                                    // see: https://casino.ru/garden-of-riches-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPcpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat,
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
