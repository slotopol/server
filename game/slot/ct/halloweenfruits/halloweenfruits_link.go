//go:build !prod || full || ct

package halloweenfruits

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed halloweenfruits_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Halloween Fruits", LNum: 30, Date: game.Date(2020, 11, 26)},      // see: https://www.slotsmate.com/software/ct-interactive/ct-gaming-halloween-fruits
		{Prov: "CT Interactive", Name: "The Power of Ramesses", LNum: 30, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/the-power-of-ramesses
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgreel |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/halloweenfruits/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/halloweenfruits/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
