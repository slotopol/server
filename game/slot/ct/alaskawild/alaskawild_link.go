//go:build !prod || full || ct

package alaskawild

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed alaskawild_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Alaska Wild", LNum: 50, Date: game.Date(2018, 12, 31)},     // see: https://www.slotsmate.com/software/ct-interactive/alaska-wild
		{Prov: "CT Interactive", Name: "Pyramid of Gold", LNum: 50, Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/pyramid-of-gold
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/alaskawild/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
