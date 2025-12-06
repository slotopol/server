//go:build !prod || full || novomatic

package beetlemania

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed beetlemania_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Beetle Mania", LNum: 9, Date: game.Year(2009)}, // see: https://www.slotsmate.com/software/novomatic/beetle-mania
		{Prov: "Novomatic", Name: "Beetle Mania Deluxe", LNum: 10, Date: game.Date(2007, 11, 13)},
		{Prov: "Novomatic", Name: "Hot Target", LNum: 9, Date: game.Date(2009, 9, 12)}, // see: https://www.slotsmate.com/software/novomatic/hot-target
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
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
	game.DataRouter["novomatic/beetlemania/bon"] = &ReelsBon
	game.DataRouter["novomatic/beetlemania/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
