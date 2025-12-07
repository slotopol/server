//go:build !prod || full || agt

package extraspin

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed extraspin_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Extra Spin", LNum: 10}, // see: https://agtsoftware.com/games/agt/extraspin
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
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
	game.DataRouter["agt/extraspin/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
