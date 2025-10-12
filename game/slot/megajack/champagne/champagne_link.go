//go:build !prod || full || megajack

package champagne

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed champagne_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Champagne", Date: game.Year(1999)},
		{Prov: "CT Interactive", Name: "Champagne and Fruits", Date: game.Date(2022, 2, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/champagne-and-fruits
		{Prov: "CT Interactive", Name: "Bloody Princess", Date: game.Date(2022, 10, 14)},    // see: https://www.slotsmate.com/software/ct-interactive/bloody-princess
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPjack |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["megajack/champagne/reel"] = &ReelsMap
	game.DataRouter["megajack/champagne/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, data)
}
