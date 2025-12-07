//go:build !prod || full || agt

package icefruits

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed icefruits_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Fruits", LNum: 20}, // see: https://agtsoftware.com/games/agt/megaice
		{Prov: "AGT", Name: "Mega Shine", LNum: 30}, // see: https://agtsoftware.com/games/agt/megashine
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/icefruits/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
