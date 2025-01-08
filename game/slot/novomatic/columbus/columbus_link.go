//go:build !prod || full || novomatic

package columbus

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Columbus"},
		{Prov: "Novomatic", Name: "Columbus Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/columbus-deluxe
		{Prov: "Novomatic", Name: "Marco Polo"},
		{Prov: "Novomatic", Name: "Holmes and Watson Deluxe"}, // see: https://www.slotsmate.com/software/novomatic/holmes-and-watson-deluxe
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgreel |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
