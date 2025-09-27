//go:build !prod || full || ct

package fortunepyramid

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fortunepyramid_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fortune Pyramid", Date: game.Date(2019, 12, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/fortune-pyramid
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
	game.DataRouter["ctinteractive/fortunepyramid/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
