//go:build !prod || full || agt

package santa

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed santa_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Santa", LNum: 10}, // see: https://agtsoftware.com/games/agt/santa
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild,
		SX: 4,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["agt/santa/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
