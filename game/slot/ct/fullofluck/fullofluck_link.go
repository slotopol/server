//go:build !prod || full || ct

package fullofluck

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fullofluck_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Full of Luck", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/full-of-luck
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPretrig |
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
	game.DataRouter["ctinteractive/fullofluck/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
