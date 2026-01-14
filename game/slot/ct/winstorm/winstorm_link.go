//go:build !prod || full || ct

package winstorm

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed winstorm_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Win Storm", LNum: 30, Date: game.Date(2020, 12, 27)},       // see: https://www.slotsmate.com/software/ct-interactive/win-storm
		{Prov: "CT Interactive", Name: "Dark Woods", LNum: 30, Date: game.Date(2020, 12, 16)},      // see: https://www.slotsmate.com/software/ct-interactive/dark-woods
		{Prov: "CT Interactive", Name: "30 Fruitata Wins", LNum: 30, Date: game.Date(2024, 2, 15)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/30-fruitata-wins
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcasc |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/winstorm/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
