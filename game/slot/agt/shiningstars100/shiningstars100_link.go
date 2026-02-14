//go:build !prod || full || agt

package shiningstars100

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed shiningstars100_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "100 Shining Stars", LNum: 100},                // see: https://agtsoftware.com/games/agt/shiningstars100
		{Prov: "AGT", Name: "50 Apples' Shine", LNum: 50},                  // see: https://agtsoftware.com/games/agt/applesshine50
		{Prov: "AGT", Name: "Red Crown", LNum: 100, Date: game.Year(2025)}, // see: https://agtsoftware.com/games/agt/redcrown
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
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
	game.DataRouter["agt/100shiningstars/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
