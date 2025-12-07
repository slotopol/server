//go:build !prod || full || agt

package aislot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed aislot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "AI", LNum: 15},          // see: https://agtsoftware.com/games/agt/aislot
		{Prov: "AGT", Name: "Tesla", LNum: 15},       // see: https://agtsoftware.com/games/agt/tesla
		{Prov: "AGT", Name: "Book of Set", LNum: 30}, // see: https://agtsoftware.com/games/agt/bookofset
		{Prov: "AGT", Name: "Pharaoh II", LNum: 30},  // see: https://agtsoftware.com/games/agt/pharaoh2
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPwsc,
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
	game.DataRouter["agt/aislot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
