//go:build !prod || full || netent

package simsalabim

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed simsalabim_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Simsalabim", Date: game.Year(2011)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
			game.GPfgreel |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["netent/simsalabim/bon"] = &ReelsBon
	game.DataRouter["netent/simsalabim/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
