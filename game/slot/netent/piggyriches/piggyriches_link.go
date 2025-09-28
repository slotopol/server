//go:build !prod || full || netent

package piggyriches

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed piggyriches_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Piggy Riches", Date: game.Date(2010, 9, 5)}, // see: https://www.slotsmate.com/software/netent/piggy-riches
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfghas |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["netent/piggyriches/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
