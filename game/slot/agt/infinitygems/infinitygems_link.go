//go:build !prod || full || agt

package infinitygems

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed infinitygems_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Infinity Gems", LNum: 20}, // see: https://agtsoftware.com/games/agt/infinitygems
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["agt/infinitygems/bon"] = &ReelsBon
	game.DataRouter["agt/infinitygems/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
