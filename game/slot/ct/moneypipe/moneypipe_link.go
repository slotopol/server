//go:build !prod || full || ct

package moneypipe

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Money Pipe", Date: game.Date(2020, 10, 1)},    // see: https://www.slotsmate.com/software/ct-interactive/money-pipe
		{Prov: "CT Interactive", Name: "Ice Rubies", Date: game.Date(2020, 12, 1)},    // see: https://www.slotsmate.com/software/ct-interactive/ice-rubies
		{Prov: "CT Interactive", Name: "More Dragons", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/more-dragons
		{Prov: "CT Interactive", Name: "Colibri Wild", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/colibri-wild
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild |
			game.GPbwild,
		SX:  5,
		SY:  4,
		SN:  len(LinePay),
		LN:  len(BetLines),
		BN:  0,
		RTP: game.MakeRtpList(ReelsMap),
	},
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
}
