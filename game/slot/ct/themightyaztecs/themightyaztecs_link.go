//go:build !prod || full || ct

package themightyaztecs

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed themightyaztecs_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "The Mighty Aztecs", Date: game.Date(2020, 1, 7)}, // see: https://www.slotsmate.com/software/ct-interactive/the-mighty-aztecs
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/themightyaztecs/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
