//go:build !prod || full || ct

package hellshot7s

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed hellshot7s_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Hell's Hot 7's", Date: game.Date(2025, 6, 30)}, // see: https://www.slotsmate.com/software/ct-interactive/hells-hot-7s
		{Prov: "CT Interactive", Name: "Hot 7's x2", Date: game.Date(2020, 12, 25)},    // see: https://www.slotsmate.com/software/ct-interactive/hot-7s-x2
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["ctinteractive/hellshot7s/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
