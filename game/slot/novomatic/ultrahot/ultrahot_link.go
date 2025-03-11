//go:build !prod || full || novomatic

package ultrahot

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Ultra Hot"},        // see: https://www.slotsmate.com/software/novomatic/ultra-hot
		{Prov: "Novomatic", Name: "Ultra Hot Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/ultrahot-deluxe
		{Prov: "Novomatic", Name: "Ultra Gems"},       // see: https://www.slotsmate.com/software/novomatic/ultra-gems
	},
	GP: game.GPsel |
		game.GPfgno |
		game.GPfill,
	SX:  3,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStat)
}
