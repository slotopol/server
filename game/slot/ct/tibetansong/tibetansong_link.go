//go:build !prod || full || ct

package tibetansong

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed tibetansong_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Tibetan Song", LNum: 25, Date: game.Date(2020, 11, 26)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/tibetan-song
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/tibetansong/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
