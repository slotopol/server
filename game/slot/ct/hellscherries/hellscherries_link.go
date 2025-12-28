//go:build !prod || full || ct

package hellscherries

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed hellscherries_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Hell's Cherries", LNum: 5, Date: game.Date(2025, 9, 30)}, // see: https://www.slotsmate.com/software/ct-interactive/hells-cherries
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgno |
			game.GPwild,
		SX: 3,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/hellscherries/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
