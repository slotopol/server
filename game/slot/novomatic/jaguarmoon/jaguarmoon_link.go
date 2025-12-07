//go:build !prod || full || novomatic

package jaguarmoon

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed jaguarmoon_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Jaguar Moon", Date: game.Date(2018, 7, 11)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		WN: 243,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["novomatic/jaguarmoon/bon"] = &ReelsBon
	game.DataRouter["novomatic/jaguarmoon/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
