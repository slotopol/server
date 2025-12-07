//go:build !prod || full || ct

package mistyforest

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed mistyforest_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Misty Forest", LNum: 20, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/misty-forest
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgmult |
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
	game.DataRouter["ctinteractive/mistyforest/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
