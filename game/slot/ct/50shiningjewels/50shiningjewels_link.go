//go:build !prod || full || ct

package fiftyshiningjewels

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed 50shiningjewels_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "50 Shining Jewels", LNum: 50, Date: game.Date(2019, 3, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/50-shining-jewels
		{Prov: "CT Interactive", Name: "50 Treasures", LNum: 50, Date: game.Date(2023, 5, 12)},      // see: https://www.livebet2.com/casino/slots/ct-interactive/50-treasures
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/50shiningjewels/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
