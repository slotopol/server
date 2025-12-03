//go:build !prod || full || novomatic

package fairyqueen

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fairyqueen_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Fairy Queen", Date: game.Date(2013, 11, 28)}, // see: https://www.slotsmate.com/software/novomatic/fairy-queen
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
	game.DataRouter["novomatic/fairyqueen/es"] = &ReelExpSym
	game.DataRouter["novomatic/fairyqueen/bon"] = &ReelsBon
	game.DataRouter["novomatic/fairyqueen/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
