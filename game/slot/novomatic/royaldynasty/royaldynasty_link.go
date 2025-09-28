//go:build !prod || full || novomatic

package royaldynasty

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed royaldynasty_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Royal Dynasty", Date: game.Date(2013, 4, 1)}, // see: https://www.slotsmate.com/software/novomatic/royal-dynasty
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgreel |
			game.GPscat |
			game.GPwild |
			game.GPwmult |
			game.GPwturn,
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
	game.DataRouter["novomatic/royaldynasty/bon"] = &ReelsBon
	game.DataRouter["novomatic/royaldynasty/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
