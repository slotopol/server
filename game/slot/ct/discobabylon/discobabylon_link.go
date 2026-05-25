//go:build !prod || full || ct

package discobabylon

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed discobabylon_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Disco Babylon", LNum: 10, Date: game.Date(2013, 6, 15)}, // see: https://www.livebet.com/casino/slots/ct-interactive/disco-babylon
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["ctinteractive/discobabylon/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
