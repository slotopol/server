//go:build !prod || full || novomatic

package fruitilicious

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitilicious_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Fruitilicious", LNum: 5, Date: game.Date(2009, 1, 14)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPrpay |
			game.GPfgno,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/fruitilicious/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
