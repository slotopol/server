//go:build !prod || full || ct

package sugartown

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed sugartown_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Sugar Town", Date: game.Date(2024, 11, 30)},        // see: https://www.slotsmate.com/software/ct-interactive/sugar-town
		{Prov: "CT Interactive", Name: "Guardian of Asgard", Date: game.Date(2024, 6, 15)}, // see: https://www.slotsmate.com/software/ct-interactive/guardian-of-asgard
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPcpay |
			game.GPcasc |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: 0,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/sugartown/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
