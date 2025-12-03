//go:build !prod || full || ct

package mightyrex

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed mightyrex_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Mighty Rex", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/mighty-rex
		{Prov: "CT Interactive", Name: "Bavarian Forest", Date: game.Year(2019)},    // see: https://www.slotsmate.com/software/ct-interactive/bavarian-forest
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/mightyrex/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/mightyrex/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
