//go:build !prod || full || netent

package tikiwonders

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed tikiwonders_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Tiki Wonders", Date: game.Year(2008)},
		{Prov: "NetEnt", Name: "Geisha Wonders", Date: game.Year(2013)},
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
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["netent/tikiwonders/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
