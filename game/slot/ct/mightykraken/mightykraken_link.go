//go:build !prod || full || ct

package mightykraken

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed mightykraken_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Mighty Kraken", LNum: 5, Date: game.Date(2020, 11, 25)},    // see: https://www.slotsmate.com/software/ct-interactive/mighty-kraken
		{Prov: "CT Interactive", Name: "FC Magic", LNum: 5, Date: game.Date(2024, 8, 31)},          // see: https://www.slotsmate.com/software/ct-interactive/fc-magic
		{Prov: "CT Interactive", Name: "Fruits of Desire", LNum: 5, Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/fruits-of-desire
		{Prov: "CT Interactive", Name: "Lucky Clover", LNum: 5, Date: game.Date(2016, 6, 30)},      // see: https://www.slotsmate.com/software/ct-interactive/casino-technology-lucky-clover
		{Prov: "CT Interactive", Name: "Wonder 7's", LNum: 5, Date: game.Date(2025, 7, 31)},        // see: https://www.slotsmate.com/software/ct-interactive/wonder-7s
		{Prov: "CT Interactive", Name: "Burning Flower", LNum: 5, Date: game.Date(2025, 11, 16)},   // see: https://www.slotsmate.com/software/ct-interactive/burning-flower
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/mightykraken/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
