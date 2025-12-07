//go:build !prod || full || aristocrat

package redroo

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed redroo_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Aristocrat", Name: "Redroo", Date: game.Date(2017, 6, 28)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		WN: 1024,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["aristocrat/redroo/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
