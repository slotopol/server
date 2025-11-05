//go:build !prod || full || ct

package purplefruits

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed purplefruits_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Purple Fruits", Date: game.Date(2020, 12, 21)}, // see: https://www.slotsmate.com/software/ct-interactive/purple-fruits
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat,
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
	game.DataRouter["ctinteractive/purplefruits/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
