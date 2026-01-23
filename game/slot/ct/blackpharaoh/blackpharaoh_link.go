//go:build !prod || full || ct

package blackpharaoh

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed blackpharaoh_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Black Pharaoh", LNum: 20, Date: game.Date(2019, 1, 1)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/black-pharaoh
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
	game.DataRouter["ctinteractive/blackpharaoh/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
