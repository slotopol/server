//go:build !prod || full || netent

package groovysixties

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed groovysixties_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Groovy Sixties", LNum: 40, Date: game.Date(2009, 6, 15)},
		{Prov: "NetEnt", Name: "Funky Seventies", LNum: 40, Date: game.Date(2009, 6, 15)}, // See: https://www.youtube.com/watch?v=a-qF9ZOpRP0
		{Prov: "NetEnt", Name: "Super Eighties", LNum: 40, Date: game.Date(2009, 6, 15)},  // See: https://www.youtube.com/watch?v=Wj49gwfRtz8
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
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
	game.DataRouter["netent/groovysixties/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
