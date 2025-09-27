//go:build !prod || full || novomatic

package oliversbar

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed oliversbar_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Oliver's Bar", Date: game.Year(2001)}, // see: https://casino.ru/olivers-bar-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["novomatic/oliversbar/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
