//go:build !prod || full || novomatic

package jewels

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed jewels_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Jewels", Date: game.Date(2007, 3, 1)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPcpay |
			game.GPlsel |
			game.GPfgno,
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
	game.DataRouter["novomatic/jewels/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
