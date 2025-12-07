//go:build !prod || full || agt

package fruitqueen

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitqueen_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Fruit Queen", LNum: 18}, // see: https://agtsoftware.com/games/agt/fruitqueen
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 6,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/fruitqueen/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
