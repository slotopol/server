//go:build !prod || full || agt

package extraspin2

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed extraspin2_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Extra Spin II"}, // see: https://agtsoftware.com/games/agt/extraspin2
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["agt/extraspinii/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
