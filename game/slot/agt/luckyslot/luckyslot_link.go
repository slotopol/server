//go:build !prod || full || agt

package luckyslot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed luckyslot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Lucky Slot", LNum: 10}, // see: https://agtsoftware.com/games/agt/luckyslot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/luckyslot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
