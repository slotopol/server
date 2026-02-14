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
		{Prov: "CT Interactive", Name: "Cherry Crown", LNum: 20, Date: game.Date(2020, 7, 1)},     // see: https://www.slotsmate.com/software/ct-interactive/cherry-crown
		{Prov: "CT Interactive", Name: "Satyr and Nymph", LNum: 20, Date: game.Date(2020, 11, 1)}, // see: https://www.slotsmate.com/software/ct-interactive/satyr-and-nymph
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPewild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/cherrycrown/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
