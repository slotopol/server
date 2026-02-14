//go:build !prod || full || agt

package icefruits6x3

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed icefruits6x3_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Fruits 6 reels", LNum: 20}, // see: https://agtsoftware.com/games/agt/6megaice
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 6,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/icefruits6reels/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
