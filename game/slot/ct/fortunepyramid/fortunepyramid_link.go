//go:build !prod || full || ct

package fortunepyramid

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fortunepyramid_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fortune Pyramid", LNum: 20, Date: game.Date(2019, 12, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/fortune-pyramid
		{Prov: "CT Interactive", Name: "Fortune Fish", LNum: 25, Date: game.Date(2020, 11, 26)},   // see: https://www.slotsmate.com/software/ct-interactive/fortune-fish
		{Prov: "CT Interactive", Name: "Clover Gems", LNum: 25, Date: game.Date(2021, 8, 2)},      // see: https://www.livebet2.com/casino/slots/ct-interactive/clover-gems
		{Prov: "CT Interactive", Name: "Banana Party", LNum: 25, Date: game.Date(2021, 8, 2)},     // see: https://www.livebet2.com/casino/slots/ct-interactive/banana-party
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	game.DataRouter["ctinteractive/fortunepyramid/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
