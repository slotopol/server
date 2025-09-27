//go:build !prod || full || igt

package richgirl

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed richgirl_bon.yaml
var rbon []byte

//go:embed richgirl_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Rich Girl", Date: game.Date(2014, 5, 13)}, // see: https://www.slotsmate.com/software/igt/rich-girl
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
		SN: len(LinePayReg),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["igt/richgirl/bon"] = &ReelsBon
	game.DataRouter["igt/richgirl/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, rbon, reels)
}
