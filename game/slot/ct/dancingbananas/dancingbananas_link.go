//go:build !prod || full || ct

package dancingbananas

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed dancingbananas_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Dancing Bananas", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/dancing-bananas
		{Prov: "CT Interactive", Name: "Dancing Dragons", Date: game.Date(2020, 12, 11)}, // see: https://www.slotsmate.com/software/ct-interactive/dancing-dragons
		{Prov: "CT Interactive", Name: "Jaguar Warrior", Date: game.Date(2020, 11, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/jaguar-warrior
		{Prov: "CT Interactive", Name: "Clover Joker", Date: game.Date(2021, 8, 8)},      // see: https://www.slotsmate.com/software/ct-interactive/clover-joker
		{Prov: "CT Interactive", Name: "Lord of Fortune", Date: game.Date(2022, 1, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/lord-of-fortune
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
	game.DataRouter["ctinteractive/dancingbananas/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
