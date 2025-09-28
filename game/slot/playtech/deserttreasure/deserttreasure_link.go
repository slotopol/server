//go:build !prod || full || playtech

package deserttreasure

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed deserttreasure_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Desert Treasure", Date: game.Date(2005, 11, 1)}, // see: https://www.slotsmate.com/software/playtech/playtech-desert-treasure
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
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
	game.DataRouter["playtech/deserttreasure/bon"] = &ReelsBon
	game.DataRouter["playtech/deserttreasure/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
