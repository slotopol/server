//go:build !prod || full || ct

package fruitgalaxy

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitgalaxy_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fruit Galaxy", LNum: 40, Date: game.Date(2021, 5, 5)},     // see: https://www.slotsmate.com/software/ct-interactive/fruit-galaxy
		{Prov: "CT Interactive", Name: "Lord of Luck", LNum: 40, Date: game.Date(2020, 12, 18)},   // see: https://www.slotsmate.com/software/ct-interactive/lord-of-luck
		{Prov: "CT Interactive", Name: "Queen of Flames", LNum: 40, Date: game.Date(2020, 11, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/queen-of-flames
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/fruitgalaxy/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
