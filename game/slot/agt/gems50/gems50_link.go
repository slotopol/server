//go:build !prod || full || agt

package gems50

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed gems50_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "50 Gems"}, // see: https://agtsoftware.com/games/agt/gems50
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
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
	game.DataRouter["agt/50gems/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
