//go:build !prod || full || netent

package thrillspin

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed thrillspin_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Thrill Spin", LNum: 15, Date: game.Year(2009)},
		{Prov: "NetEnt", Name: "Viking's Treasure", LNum: 15, Date: game.Year(2006)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["netent/thrillspin/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
