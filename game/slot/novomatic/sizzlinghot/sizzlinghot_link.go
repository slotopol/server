//go:build !prod || full || novomatic

package sizzlinghot

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Sizzling Hot", Date: game.Date(2003, 3, 6)},          // see: https://www.slotsmate.com/software/novomatic/novomatic-sizzling-hot
		{Prov: "Novomatic", Name: "Sizzling Hot Deluxe", Date: game.Date(2007, 11, 13)}, // see: https://www.slotsmate.com/software/novomatic/sizzling-hot-deluxe
		{Prov: "Novomatic", Name: "Age of Heroes", Date: game.Date(2016, 4, 15)},        // see: https://www.slotsmate.com/software/novomatic/age-of-heroes
		{Prov: "Novomatic", Name: "Hot Cubes", Date: game.Date(2007, 6, 15)},            // see: https://www.slotsmate.com/software/novomatic/hot-cubes
		{Prov: "Novomatic", Name: "Diamond 7", Date: game.Date(2012, 11, 15)},           // see: https://www.slotsmate.com/software/novomatic/diamond-7
		{Prov: "Novomatic", Name: "Fruits 'n Royals", Date: game.Date(2009, 12, 29)},    // see: https://www.slotsmate.com/software/novomatic/fruits-n-royals
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
