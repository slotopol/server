//go:build !prod || full || ct

package neonbananas

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed neonbananas_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Neon Bananas", LNum: 20, Date: game.Date(2021, 6, 6)},  // see: https://www.slotsmate.com/software/ct-interactive/neon-bananas
		{Prov: "CT Interactive", Name: "Mighty Moon", LNum: 20, Date: game.Date(2021, 7, 7)},   // see: https://www.slotsmate.com/software/ct-interactive/mighty-moon
		{Prov: "CT Interactive", Name: "Clover Wheel", LNum: 20, Date: game.Date(2020, 12, 4)}, // see: https://www.slotsmate.com/software/ct-interactive/clover-wheel
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 2,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/neonbananas/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
