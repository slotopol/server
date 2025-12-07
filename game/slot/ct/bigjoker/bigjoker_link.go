//go:build !prod || full || ct

package bigjoker

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed bigjoker_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Big Joker", LNum: 20, Date: game.Date(2018, 12, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/big-joker
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPrwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/bigjoker/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
