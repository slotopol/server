//go:build !prod || full || megajack

package aztecgold

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed aztecgold_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Megajack", Name: "Aztec Gold", LNum: 21, Date: game.Year(1999)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPjack |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["megajack/aztecgold/reel"] = &ReelsMap
	game.DataRouter["megajack/aztecgold/jack"] = &JackMap
	game.LoadMap = append(game.LoadMap, data)
}
