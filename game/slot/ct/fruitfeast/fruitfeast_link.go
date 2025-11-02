//go:build !prod || full || ct

package fruitfeast

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitfeast_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fruit Feast", Date: game.Date(2020, 9, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/fruit-feast
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/fruitfeast/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
