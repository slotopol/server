//go:build !prod || full || ct

package vikingsfun

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed vikingsfun_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Vikings Fun", Date: game.Date(2020, 10, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/vikings-fun
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPretrig |
			game.GPscat |
			game.GPwild |
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
	game.DataRouter["ctinteractive/vikingsfun/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
