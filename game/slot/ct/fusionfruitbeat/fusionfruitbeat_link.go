//go:build !prod || full || ct

package fusionfruitbeat

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fusionfruitbeat_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fusion Fruit Beat", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/fusion-fruit-beat
		{Prov: "CT Interactive", Name: "Lady Emerald", Date: game.Date(2021, 3, 1)},        // see: https://www.slotsmate.com/software/ct-interactive/lady-emerald
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/fusionfruitbeat/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
