//go:build !prod || full || ct

package cherrycrown

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed cherrycrown_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Cherry Crown", Date: game.Date(2020, 7, 1)},     // see: https://www.slotsmate.com/software/ct-interactive/cherry-crown
		{Prov: "CT Interactive", Name: "Satyr and Nymph", Date: game.Date(2020, 11, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/satyr-and-nymph
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild |
			game.GPbwild,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/cherrycrown/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
