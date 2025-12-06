//go:build !prod || full || novomatic

package columbus

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed columbus_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Columbus", LNum: 9, Date: game.Year(2005)},                         // see: https://casino.ru/columbus-novomatic/
		{Prov: "Novomatic", Name: "Columbus Deluxe", LNum: 10, Date: game.Date(2008, 3, 19)},          // see: https://www.slotsmate.com/software/novomatic/columbus-deluxe
		{Prov: "Novomatic", Name: "Marco Polo", LNum: 9, Date: game.Year(2008)},                       // see: https://casino.ru/marco-polo-novomatic/
		{Prov: "Novomatic", Name: "Holmes and Watson Deluxe", LNum: 10, Date: game.Date(2018, 3, 15)}, // see: https://www.slotsmate.com/software/novomatic/holmes-and-watson-deluxe
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["novomatic/columbus/bon"] = &ReelsBon
	game.DataRouter["novomatic/columbus/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
