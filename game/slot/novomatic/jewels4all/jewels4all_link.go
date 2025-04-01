//go:build !prod || full || novomatic

package jewels4all

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Jewels 4 All"},
	},
	GP: game.GPcpay |
		game.GPsel |
		game.GPfgno |
		game.GPbwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ChanceMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
