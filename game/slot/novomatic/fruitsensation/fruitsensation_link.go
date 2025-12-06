//go:build !prod || full || novomatic

package fruitsensation

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitsensation_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Fruit Sensation", LNum: 10, Date: game.Date(2012, 6, 15)}, // see: https://casino.ru/fruit-sensation-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/fruitsensation/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
