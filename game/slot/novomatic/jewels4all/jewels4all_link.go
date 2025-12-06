//go:build !prod || full || novomatic

package jewels4all

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed jewels4all_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Jewels 4 All", LNum: 9, Date: game.Year(2009)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPapay |
			game.GPlsel |
			game.GPfgno |
			game.GPbwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ChanceMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/jewels4all/reel"] = &Reels
	game.DataRouter["novomatic/jewels4all/chance"] = &ChanceMap
	game.LoadMap = append(game.LoadMap, data)
}
