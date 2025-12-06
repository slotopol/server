//go:build !prod || full || ct

package groovyautomat

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed groovyautomat_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Groovy Automat", LNum: 5, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/groovy-automat
		{Prov: "CT Interactive", Name: "Golden Amulet", LNum: 5, Date: game.Date(2020, 11, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/golden-amulet
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgno |
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
	game.DataRouter["ctinteractive/groovyautomat/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
