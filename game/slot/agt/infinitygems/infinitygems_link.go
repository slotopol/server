//go:build !prod || full || agt

package infinitygems

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed infinitygems_bon.yaml
var rbon []byte

//go:embed infinitygems_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Infinity Gems"},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["agt/infinitygems/bon"] = &ReelsBon
	game.DataRouter["agt/infinitygems/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, rbon, reels)
}
