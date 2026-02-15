//go:build !prod || full || novomatic

package oliversbar

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed oliversbar_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Oliver's Bar", LNum: 9, Date: game.Year(2001)}, // see: https://casino.ru/olivers-bar-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["novomatic/oliversbar/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
