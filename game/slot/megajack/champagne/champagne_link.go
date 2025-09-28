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
