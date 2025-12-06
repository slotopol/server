//go:build !prod || full || agt

package panda

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed panda_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Panda", LNum: 27}, // see: https://agtsoftware.com/games/agt/panda
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgseq |
			game.GPscat |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["agt/panda/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
