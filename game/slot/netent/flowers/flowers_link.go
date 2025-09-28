//go:build !prod || full || netent

package flowers

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed flowers_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Flowers", Date: game.Date(2013, 11, 11)}, // see: https://games.netent.com/video-slots/flowers/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPfgmult |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["netent/flowers/bon"] = &ReelsBon
	game.DataRouter["netent/flowers/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
