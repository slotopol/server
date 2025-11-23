//go:build !prod || full || ct

package kingofclovers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed kingofclovers_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "King of Clovers", Date: game.Date(2025, 3, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/king-of-clovers
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfill |
			game.GPcasc |
			game.GPretrig |
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
	game.DataRouter["ctinteractive/kingofclovers/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/kingofclovers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
