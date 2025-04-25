//go:build !prod || full || megajack

package aztecgold

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Aztec Gold", Year: 1988},
	},
	GP: game.GPlpay |
		game.GPsel |
		game.GPjack |
		game.GPfgno |
		game.GPscat |
		game.GPrwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  1,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
