//go:build !prod || full || novomatic

package plentyontwenty

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed plentyontwenty_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Plenty on Twenty", Date: game.Year(2006)}, // see: https://casino.ru/plenty-on-twenty-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgno |
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
	game.DataRouter["novomatic/plentyontwenty/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
