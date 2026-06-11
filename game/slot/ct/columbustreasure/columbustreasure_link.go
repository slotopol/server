//go:build !prod || full || ct

package columbustreasure

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed columbustreasure_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Columbus Treasure", LNum: 10, Date: game.Date(2017, 10, 5)}, // see: https://www.livebet.com/casino/slots/ct-interactive/columbus-treasure
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/columbustreasure/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
