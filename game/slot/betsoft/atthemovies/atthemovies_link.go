//go:build !prod || full || betsoft

package atthemovies

import (
	"github.com/slotopol/server/game"
)

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "BetSoft", Name: "At the Movies", Date: game.Year(2012)}, // see: https://www.slotsmate.com/software/betsoft/at-the-movies
		{Prov: "BetSoft", Name: "Sushi Bar", Date: game.Year(2014)},     // see: https://www.slotsmate.com/software/betsoft/sushi-bar
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
}
