//go:build !prod || full || agt

package sevenhot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed sevenhot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Seven Hot", LNum: 20},                          // see: https://agtsoftware.com/games/agt/sevenhot20
		{Prov: "AGT", Name: "Live Fruits", LNum: 20, Date: game.Year(2025)}, // see: https://agtsoftware.com/games/agt/livefruits
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/sevenhot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
