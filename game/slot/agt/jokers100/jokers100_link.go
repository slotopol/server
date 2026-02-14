//go:build !prod || full || agt

package jokers100

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed jokers100_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Jokers", LNum: 100},    // see: https://agtsoftware.com/games/agt/jokers100
		{Prov: "AGT", Name: "50 Happy Santa", LNum: 50}, // see: https://agtsoftware.com/games/agt/happysanta50
		{Prov: "AGT", Name: "40 Bigfoot", LNum: 40},     // see: https://agtsoftware.com/games/agt/bigfoot40
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/100jokers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
