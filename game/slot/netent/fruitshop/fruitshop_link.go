//go:build !prod || full || netent

package fruitshop

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitshop_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Fruit Shop", LNum: 15, Date: game.Date(2011, 9, 15)}, // see: https://games.netent.com/video-slots/fruit-shop/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgmult |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["netent/fruitshop/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
