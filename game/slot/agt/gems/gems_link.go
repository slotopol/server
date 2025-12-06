//go:build !prod || full || agt

package gems

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed gems_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Gems", LNum: 20}, // see: https://agtsoftware.com/games/agt/gems20
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["agt/gems/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
