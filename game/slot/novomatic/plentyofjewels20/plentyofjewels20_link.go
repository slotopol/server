//go:build !prod || full || novomatic

package plentyofjewels20

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed plentyofjewels20_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Plenty of Jewels 20 hot", LNum: 20, Date: game.Date(2015, 11, 1)}, // see: https://www.slotsmate.com/software/novomatic/plenty-of-jewels-20-hot
		{Prov: "Novomatic", Name: "Plenty of Fruit 20 hot", LNum: 20, Date: game.Date(2015, 11, 4)},  // see: https://www.slotsmate.com/software/novomatic/plenty-of-fruit-20-hot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/plentyofjewels20hot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
