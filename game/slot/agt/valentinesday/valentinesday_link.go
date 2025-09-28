//go:build !prod || full || agt

package valentinesday

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed valentinesday_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Valentine's Day"},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfill |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["agt/valentinesday/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
