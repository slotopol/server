//go:build !prod || full || ct

package groovypowers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed groovypowers_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Groovy Powers", LNum: 20, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/groovy-powers
		{Prov: "CT Interactive", Name: "Space Fruits", LNum: 20, Date: game.Date(2020, 11, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/space-fruits
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPbmode |
			game.GPscat |
			game.GPwild |
			game.GPrwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["ctinteractive/groovypowers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
