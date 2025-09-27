//go:build !prod || full || agt

package shiningstars

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed shiningstars_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Shining Stars"},
		{Prov: "AGT", Name: "Green Hot"},     // see: https://agtsoftware.com/games/agt/greenhot
		{Prov: "AGT", Name: "Apples' Shine"}, // see: https://agtsoftware.com/games/agt/applesshine
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
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
	game.DataRouter["agt/shiningstars/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
