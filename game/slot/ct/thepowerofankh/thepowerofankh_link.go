//go:build !prod || full || ct

package thepowerofankh

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed thepowerofankh_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "The Power of Ankh", LNum: 25, Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/the-power-of-ankh
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgtwic |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["ctinteractive/thepowerofankh/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/thepowerofankh/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
