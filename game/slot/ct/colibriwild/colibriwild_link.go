//go:build !prod || full || ct

package colibriwild

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed colibriwild_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Colibri Wild", LNum: 40, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/colibri-wild
		{Prov: "CT Interactive", Name: "Ice Rubies", LNum: 40, Date: game.Date(2020, 12, 1)},    // see: https://www.slotsmate.com/software/ct-interactive/ice-rubies
		{Prov: "CT Interactive", Name: "Fire Dozen", LNum: 40, Date: game.Date(2020, 11, 26)},   // see: https://www.slotsmate.com/software/ct-interactive/fire-dozen
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPewild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/colibriwild/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
