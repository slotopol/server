//go:build !prod || full || novomatic

package dragonsdeep

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed dragonsdeep_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dragon's Deep", LNum: 25, Date: game.Date(2015, 12, 16)}, // see: https://www.slotsmate.com/software/novomatic/dragons-deep
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPscat |
			game.GPwild |
			game.GPwturn,
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
	game.DataRouter["novomatic/dragonsdeep/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
