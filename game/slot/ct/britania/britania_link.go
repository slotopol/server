//go:build !prod || full || ct

package britania

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed britania_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Britania", LNum: 20, Date: game.Date(2020, 10, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/britania
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/britania/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
