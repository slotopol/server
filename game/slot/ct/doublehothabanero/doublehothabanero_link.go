//go:build !prod || full || ct

package doublehothabanero

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed doublehothabanero_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Double Hot Habanero", LNum: 5, Date: game.Date(2020, 12, 18)}, // see: https://www.slotsmate.com/software/ct-interactive/double-hot-habanero
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	game.DataRouter["ctinteractive/doublehothabanero/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
