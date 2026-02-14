//go:build !prod || full || agt

package suncity

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed suncity_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Sun City", LNum: 30, Date: game.Year(2024)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
			game.GPfgreel |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["agt/suncity/bon"] = &ReelsBon
	game.DataRouter["agt/suncity/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
