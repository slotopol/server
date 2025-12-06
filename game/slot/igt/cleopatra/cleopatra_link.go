//go:build !prod || full || igt

package cleopatra

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed cleopatra_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Cleopatra", LNum: 20, Date: game.Date(2012, 4, 1)}, // see: https://www.slotsmate.com/software/igt/igt-cleopatra
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["igt/cleopatra/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
