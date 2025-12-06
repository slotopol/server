//go:build !prod || full || novomatic

package megajoker

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed megajoker_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Mega Joker", LNum: 40, Date: game.Date(2013, 2, 12)}, // see: https://www.slotsmate.com/software/novomatic/mega-joker
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
	game.DataRouter["novomatic/megajoker/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
