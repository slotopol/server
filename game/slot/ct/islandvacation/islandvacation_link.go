//go:build !prod || full || ct

package islandvacation

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed islandvacation_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Island Vacation", LNum: 15, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/island-vacation
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["ctinteractive/islandvacation/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
