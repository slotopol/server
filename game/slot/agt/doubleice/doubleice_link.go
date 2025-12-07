//go:build !prod || full || agt

package doubleice

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed doubleice_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Double Ice", LNum: 27}, // see: https://agtsoftware.com/games/agt/doubleice
		{Prov: "AGT", Name: "Double Hot", LNum: 27}, // see: https://agtsoftware.com/games/agt/double
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfill |
			game.GPfgno,
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
	game.DataRouter["agt/doubleice/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
