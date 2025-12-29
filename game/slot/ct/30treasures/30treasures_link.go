//go:build !prod || full || ct

package thirtytreasures

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed 30treasures_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "30 Treasures", LNum: 30, Date: game.Date(2019, 12, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/30-treasures
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/30treasures/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
