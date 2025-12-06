//go:build !prod || full || novomatic

package inferno

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed inferno_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Inferno", LNum: 5, Date: game.Date(2014, 11, 15)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat,
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
	game.DataRouter["novomatic/inferno/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
