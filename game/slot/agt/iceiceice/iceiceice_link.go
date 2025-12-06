//go:build !prod || full || agt

package iceiceice

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed iceiceice_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Ice Ice", LNum: 10},  // see: https://agtsoftware.com/games/agt/iceiceice
		{Prov: "AGT", Name: "5 Hot Hot Hot", LNum: 5}, // see: https://agtsoftware.com/games/agt/hothothot5
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["agt/iceiceice/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
