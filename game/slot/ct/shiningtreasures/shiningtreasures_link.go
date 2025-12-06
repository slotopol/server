//go:build !prod || full || ct

package shiningtreasures

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed shiningtreasures_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Shining Treasures", LNum: 15, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/shining-treasures
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/shiningtreasures/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
