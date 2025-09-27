//go:build !prod || full || ct

package winstorm

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed winstorm_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Win Storm", Date: game.Date(2020, 12, 27)},  // see: https://www.slotsmate.com/software/ct-interactive/win-storm
		{Prov: "CT Interactive", Name: "Dark Woods", Date: game.Date(2020, 12, 16)}, // see: https://www.slotsmate.com/software/ct-interactive/dark-woods
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/winstorm/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
