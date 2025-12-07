//go:build !prod || full || novomatic

package dynastyofra

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed dynastyofra_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dynasty of Ra", LNum: 10, Date: game.Date(2014, 1, 29)}, // see: https://www.slotsmate.com/software/novomatic/dynasty-of-ra
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPwsc,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["novomatic/dynastyofra/bon"] = &ReelsBon
	game.DataRouter["novomatic/dynastyofra/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
