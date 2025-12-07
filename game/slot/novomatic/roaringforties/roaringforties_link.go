//go:build !prod || full || novomatic

package roaringforties

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed roaringforties_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Roaring Forties", LNum: 40, Date: game.Date(2013, 4, 29)},
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/roaringforties/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
