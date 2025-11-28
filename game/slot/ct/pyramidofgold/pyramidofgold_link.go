//go:build !prod || full || ct

package pyramidofgold

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed pyramidofgold_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Pyramid of Gold", Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/pyramid-of-gold
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPretrig |
			game.GPfgreel |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/pyramidofgold/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/pyramidofgold/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
