//go:build !prod || full || ct

package caveofclovers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed caveofclovers_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Cave of Clovers", LNum: 100, Date: game.Date(2025, 11, 30)}, // see: https://www.slotsmate.com/software/ct-interactive/cave-of-clovers
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/caveofclovers/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
